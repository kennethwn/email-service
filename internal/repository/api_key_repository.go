package repository

import (
	"context"
	"errors"
	"fmt"
	"worker-service/internal/dto"

	"gorm.io/gorm"
)

type apiKeyRepository struct {
	db *gorm.DB
}

type ApiKeyRepository interface {
	FetchOne(ctx context.Context, query Query) (dto.ApiKey, error)
}

func NewApiKeyRepository(db *gorm.DB) ApiKeyRepository {
	return &apiKeyRepository{
		db: db,
	}
}

func (r *apiKeyRepository) FetchOne(ctx context.Context, query Query) (dto.ApiKey, error) {
	var apiKey dto.ApiKey
	db := r.db.Model(dto.ApiKey{}).WithContext(ctx)
	db = QueryHelperDB(db, query)

	err := db.First(&apiKey).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ApiKey{}, fmt.Errorf("data record is not found")
		}
		return dto.ApiKey{}, err
	}

	return apiKey, nil
}
