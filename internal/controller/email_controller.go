package controller

import (
	"strconv"
	"worker-service/internal/dto"
	"worker-service/internal/pkg/error_wrap"
	"worker-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

const (
	EmailPath         = "/emails"
	EmailByIdPath     = "/emails/:id"
	EmailSendBulkPath = "/emails/bulk"
	EmailRetryPath    = "/emails/retry"
)

type emailController struct {
	emailUsecase usecase.EmailUsecase
}

type EmailController interface {
	ListEmail(ctx *gin.Context)
	ListEmailByID(ctx *gin.Context)
	RetryEmail(ctx *gin.Context)
	SendEmail(ctx *gin.Context)
}

func NewEmailController(emailUsecase usecase.EmailUsecase) EmailController {
	return &emailController{
		emailUsecase: emailUsecase,
	}
}

func (c *emailController) ListEmail(ctx *gin.Context) {
	pagination := ParsePagination(ctx)

	isAscendingStr := ctx.Query("is_ascending")
	isAscending, err := strconv.ParseBool(isAscendingStr)
	if err != nil {
		isAscending = false
	}

	status := ctx.QueryArray("status")

	data, err := c.emailUsecase.ListEmail(ctx, usecase.ListEmailRequestQuery{
		Page:        pagination.Page,
		Limit:       pagination.Limit,
		IsAscending: isAscending,
		Status:      status,
	})
	if err != nil {
		dto.WriteErrorResponseJSON(ctx, err)
		return
	}

	dto.SuccessResponse.Data = data
	dto.WriteResponseJSON(ctx, dto.SuccessResponse)
}

func (c *emailController) ListEmailByID(ctx *gin.Context) {
	pagination := ParsePagination(ctx)
	emailID := ctx.Param("id")

	data, err := c.emailUsecase.ListEmail(ctx, usecase.ListEmailRequestQuery{
		Page:  pagination.Page,
		Limit: pagination.Limit,
		ID:    emailID,
	})
	if err != nil {
		dto.WriteErrorResponseJSON(ctx, err)
		return
	}

	dto.SuccessResponse.Data = data
	dto.WriteResponseJSON(ctx, dto.SuccessResponse)
}

func (c *emailController) RetryEmail(ctx *gin.Context) {
	var request dto.RetryEmailRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		dto.WriteErrorResponseJSON(ctx, error_wrap.ErrBadRequest)
		return
	}

	if err := c.emailUsecase.RetryEmail(ctx, request.ID); err != nil {
		dto.WriteErrorResponseJSON(ctx, error_wrap.ErrInternalServerError)
		return
	}

	dto.SuccessResponse.Data = nil
	dto.WriteResponseJSON(ctx, dto.SuccessResponse)
}

func (c *emailController) SendEmail(ctx *gin.Context) {
	var request []dto.EmailTask
	if err := ctx.ShouldBindJSON(&request); err != nil {
		dto.WriteErrorResponseJSON(ctx, error_wrap.ErrBadRequest)
		return
	}

	resp, err := c.emailUsecase.SendEmail(ctx, request)
	if err != nil {
		dto.WriteErrorResponseJSON(ctx, error_wrap.ErrInternalServerError)
		return
	}

	dto.SuccessResponse.Data = resp
	dto.WriteResponseJSON(ctx, dto.SuccessResponse)
}
