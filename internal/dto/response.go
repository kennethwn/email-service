package dto

import (
	"errors"
	"net/http"
	"worker-service/internal/pkg/error_wrap"

	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

var SuccessResponse = BaseResponse{
	Status:  http.StatusOK,
	Message: "Success",
}

func WriteResponseJSON(ctx *gin.Context, data interface{}) {
	code := http.StatusOK
	ctx.JSON(code, data)
	return
}

func WriteErrorResponseJSON(c *gin.Context, err error) {
	switch {
	case errors.Is(err, error_wrap.ErrBadRequest):
		c.JSON(http.StatusBadRequest, BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	case errors.Is(err, error_wrap.ErrNotFound):
		c.JSON(http.StatusNotFound, BaseResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		})
	case errors.Is(err, error_wrap.ErrTooManyRequests):
		c.JSON(http.StatusTooManyRequests, BaseResponse{
			Status:  http.StatusTooManyRequests,
			Message: err.Error(),
		})
	case errors.Is(err, error_wrap.ErrForbidden), errors.Is(err, error_wrap.ErrIPorServiceBlocked):
		c.JSON(http.StatusForbidden, BaseResponse{
			Status:  http.StatusForbidden,
			Message: err.Error(),
		})
	case errors.Is(err, error_wrap.ErrInvalidToken), errors.Is(err, error_wrap.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	case errors.Is(err, error_wrap.ErrApiKeyIsMissing), errors.Is(err, error_wrap.ErrApiKeyIsInvalid):
		c.JSON(http.StatusUnauthorized, BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	case errors.Is(err, error_wrap.ErrSqlError), errors.Is(err, error_wrap.ErrInternalServerError):
		c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "unexpected error",
		})
	}
	return
}
