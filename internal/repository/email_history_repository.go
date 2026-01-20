package repository

import (
	"context"
	"errors"
	"fmt"
	"worker-service/internal/dto"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type EmailHistoryRepository interface {
	Create(ctx context.Context, email *dto.EmailHistory) error
	Update(ctx context.Context, id string, email *dto.EmailHistory) error
	FetchOne(ctx context.Context, query Query) (dto.EmailHistory, error)
	Fetch(ctx context.Context, query Query) ([]dto.EmailHistory, error)
	Count(ctx context.Context, query Query) (int64, error)
}

type emailHistoryRepository struct {
	db *gorm.DB
}

func NewEmailHistoryRepository(db *gorm.DB) EmailHistoryRepository {
	return &emailHistoryRepository{db: db}
}

func (r *emailHistoryRepository) Create(ctx context.Context, data *dto.EmailHistory) error {
	if data.ID == "" {
		data.ID = ulid.Make().String()
	}
	return r.db.Model(dto.EmailHistory{}).WithContext(ctx).Create(data).Error
}

func (r *emailHistoryRepository) Update(ctx context.Context, id string, data *dto.EmailHistory) error {
	return r.db.Model(dto.EmailHistory{}).Where("id = ?", id).WithContext(ctx).Updates(&data).Error
}

func (r *emailHistoryRepository) FetchOne(ctx context.Context, query Query) (dto.EmailHistory, error) {
	var email dto.EmailHistory
	db := r.db.Model(dto.EmailHistory{}).WithContext(ctx)
	db = QueryHelperDB(db, query)

	err := db.First(&email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.EmailHistory{}, fmt.Errorf("data record is not found")
		}
		return dto.EmailHistory{}, err
	}

	return email, nil
}

func (r *emailHistoryRepository) Fetch(ctx context.Context, query Query) ([]dto.EmailHistory, error) {
	var emailHistories []dto.EmailHistory
	db := r.db.Model(dto.EmailHistory{}).WithContext(ctx)
	db = QueryHelperDB(db, query)

	if err := db.Find(&emailHistories).Error; err != nil {
		return nil, err
	}

	return emailHistories, nil
}

func (r *emailHistoryRepository) Count(ctx context.Context, query Query) (int64, error) {
	var count int64
	db := r.db.Model(dto.EmailHistory{}).WithContext(ctx)
	db = QueryHelperDB(db, query)

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
