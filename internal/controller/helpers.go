package controller

import (
	"strconv"
	"worker-service/internal/pkg/error_wrap"

	"github.com/gin-gonic/gin"
)

type PaginationReq struct {
	Page  int
	Limit int
}

func ParsePagination(ctx *gin.Context) PaginationReq {
	var err error
	var result PaginationReq

	result.Page, err = ParseQueryToInt(ctx, "page")
	if err != nil || result.Page == 0 {
		result.Page = 1
	}
	result.Limit, err = ParseQueryToInt(ctx, "limit")
	if err != nil || result.Limit == 0 {
		result.Limit = 10
	}

	return result
}

func ParseQueryToInt(ctx *gin.Context, query string) (int, error) {
	queryStr := ctx.Query(query)
	if queryStr == "" {
		return 0, nil
	}
	queryInt, err := strconv.ParseInt(queryStr, 10, 64)
	if err != nil {
		return 0, error_wrap.ErrBadRequest
	}

	return int(queryInt), nil
}
