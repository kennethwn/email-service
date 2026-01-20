package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"
	"worker-service/config"
	"worker-service/internal/dto"
	"worker-service/internal/pkg/error_wrap"
	"worker-service/internal/pkg/redis"
	"worker-service/internal/repository"
	"worker-service/internal/repository/unitofwork"
	"worker-service/internal/services"

	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc/pool"
)

type EmailUsecase interface {
	ListEmail(ctx context.Context, query ListEmailRequestQuery) (ListEmailResponse, error)
	RetryEmail(ctx context.Context, id string) error
	SendEmail(ctx context.Context, data []dto.EmailTask) (SendEmailResponse, error)
}

type emailUsecase struct {
	cfg              *config.AppConfig
	emailHistoryRepo repository.EmailHistoryRepository
	emailService     services.EmailService
	uow              unitofwork.UnitOfWork
	redisClient      *redis.RedisClient[dto.EmailTask]
}

type sendEmailWorkerResult struct {
	Success dto.EmailTask     `json:"success"`
	Failed  map[string]string `json:"failed"`
}

type SendEmailResponse struct {
	Success     int64               `json:"success"`
	Failed      int64               `json:"failed"`
	SuccessData []dto.EmailTask     `json:"success_data"`
	FailedData  []map[string]string `json:"failed_data"`
}

type ListEmailResponse struct {
	Header PaginationHeader   `json:"header"`
	List   []dto.EmailHistory `json:"list"`
}

type ListEmailRequestQuery struct {
	StartAt     time.Time
	EndAt       time.Time
	Page        int
	Limit       int
	IsAscending bool
	ID          string
	Status      []string
}

func NewEmailUsecase(cfg *config.AppConfig, emailHistoryRepo repository.EmailHistoryRepository, uow unitofwork.UnitOfWork, emailService services.EmailService, redisClient *redis.RedisClient[dto.EmailTask]) EmailUsecase {
	return &emailUsecase{
		cfg:              cfg,
		emailHistoryRepo: emailHistoryRepo,
		uow:              uow,
		emailService:     emailService,
		redisClient:      redisClient,
	}
}

func (u *emailUsecase) ListEmail(ctx context.Context, query ListEmailRequestQuery) (ListEmailResponse, error) {
	q, err := u.buildEmailQueryDetail(query)
	if err != nil {
		return ListEmailResponse{}, error_wrap.ErrBadRequest
	}

	totalData, err := u.emailHistoryRepo.Count(ctx, q)
	if err != nil {
		return ListEmailResponse{}, error_wrap.ErrSqlError
	}

	if totalData <= 0 {
		return ListEmailResponse{
			Header: PaginationHeader{
				CurrentPage: int64(query.Page),
				PerPage:     int64(query.Limit),
				TotalData:   0,
				TotalPages:  0,
			},
		}, nil
	}

	totalPages := math.Ceil(float64(totalData) / float64(query.Limit))
	q.Page = query.Page
	q.Limit = query.Limit
	if query.IsAscending {
		q.Order = "ASC"
	}

	data, err := u.emailHistoryRepo.Fetch(ctx, q)
	if err != nil {
		return ListEmailResponse{}, error_wrap.ErrSqlError
	}

	response := ListEmailResponse{
		Header: PaginationHeader{
			CurrentPage: int64(query.Page),
			PerPage:     int64(query.Limit),
			TotalData:   int64(len(data)),
			TotalPages:  int64(totalPages),
		},
		List: data,
	}

	return response, nil
}

func (u *emailUsecase) RetryEmail(ctx context.Context, id string) error {
	// Fetch the email
	email, err := u.emailHistoryRepo.FetchOne(ctx, repository.Query{
		Query:  "id = ? AND status = ?",
		Values: []interface{}{id, uint(dto.EmailHistoryPending)},
	})
	if err != nil {
		logrus.Error("error fetching email: ", err)
		return error_wrap.ErrSqlError
	}

	err = u.uow.Do(ctx, func(uows unitofwork.UnitOfWorkStore) error {
		// Send the email
		if err := u.emailService.SendEmail(ctx, dto.EmailTask{
			From:    email.From,
			To:      email.To,
			Subject: email.Subject,
			Body:    email.Body,
		}); err != nil {
			logrus.Error("error retry email: ", err)
			return error_wrap.ErrInternalServerError
		}

		// Update email
		email.Status = uint(dto.EmailHistorySuccess)
		email.IsActive = false
		if err := u.emailHistoryRepo.Update(ctx, id, &email); err != nil {
			logrus.Error("error updating existing email: ", err)
			return error_wrap.ErrSqlError
		}

		return nil
	}, sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		logrus.Error("error uow: ", err)
		return error_wrap.ErrInternalServerError
	}

	return nil
}

func (u *emailUsecase) SendEmail(ctx context.Context, data []dto.EmailTask) (SendEmailResponse, error) {
	var (
		resChan = make(chan sendEmailWorkerResult, len(data))
		success []dto.EmailTask
		failed  []map[string]string
	)

	pooler := pool.New().WithMaxGoroutines(10)
	for _, d := range data {
		mail := d

		pooler.Go(func() {
			mappingError := make(map[string]string)
			key := fmt.Sprintf("%s:%s", mail.To, d.Subject)
			if sendErr := u.emailService.SendEmail(ctx, mail); sendErr != nil {
				logrus.Error("error sending email: ", sendErr)
				createErr := u.emailHistoryRepo.Create(ctx, &dto.EmailHistory{
					From:     mail.From,
					To:       mail.To,
					Subject:  mail.Subject,
					Body:     mail.Body,
					Status:   uint(dto.EmailHistoryPending),
					IsActive: true,
				})
				if createErr != nil {
					logrus.Error("error creating email history: ", createErr)
					mappingError[key] = createErr.Error()
				} else {
					mappingError[key] = sendErr.Error()
				}
				resChan <- sendEmailWorkerResult{Failed: mappingError}
				return
			}

			resChan <- sendEmailWorkerResult{Success: mail}
		})
	}
	pooler.Wait()
	close(resChan)

	for res := range resChan {
		if res.Failed != nil {
			failed = append(failed, res.Failed)
		}
		if res.Success != (dto.EmailTask{}) {
			success = append(success, res.Success)
		}
	}

	successCount := int64(len(success))
	failedCount := int64(len(failed))

	var response = SendEmailResponse{
		Success:     successCount,
		Failed:      failedCount,
		SuccessData: success,
		FailedData:  failed,
	}

	return response, nil
}

func (u *emailUsecase) buildEmailQueryDetail(request ListEmailRequestQuery) (repository.Query, error) {
	query := []string{}
	emailHistoryQuery := repository.Query{
		Sort:  "updated_at",
		Order: "DESC",
	}

	var listStatus []dto.EmailHistoryStatus
	for _, item := range request.Status {
		listStatus = append(listStatus, dto.EmailHistoryStatusTypeSelector[item])
	}

	if len(listStatus) > 0 {
		query = append(query, "status IN (?)")
		emailHistoryQuery.Values = append(emailHistoryQuery.Values, listStatus)
	}

	if request.ID != "" {
		query = append(query, "id = ?")
		emailHistoryQuery.Values = append(emailHistoryQuery.Values, request.ID)
	}

	if !request.StartAt.IsZero() && !request.EndAt.IsZero() {
		if request.StartAt.After(request.EndAt) {
			return repository.Query{}, error_wrap.ErrBadRequest
		}
		emailHistoryQuery.StartAt = request.StartAt
		emailHistoryQuery.EndAt = request.EndAt
	}

	emailHistoryQuery.Query = strings.Join(query, " AND ")
	return emailHistoryQuery, nil
}
