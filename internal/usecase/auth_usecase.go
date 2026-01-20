package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"worker-service/internal/dto"
	"worker-service/internal/repository"
)

type authUsecase struct {
	apiKeyRepo repository.ApiKeyRepository
}

type AuthUsecase interface {
	VerifyAPIKey(ctx context.Context, key string) (dto.VerifyAPIKeyResponse, error)
}

func NewAuthUsecase(apiKeyRepo repository.ApiKeyRepository) AuthUsecase {
	return &authUsecase{
		apiKeyRepo: apiKeyRepo,
	}
}

func (u *authUsecase) VerifyAPIKey(ctx context.Context, key string) (dto.VerifyAPIKeyResponse, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return dto.VerifyAPIKeyResponse{}, err
	}

	hash := sha256.Sum256([]byte(decodedKey))
	keyHash := hex.EncodeToString(hash[:])

	apiKey, err := u.apiKeyRepo.FetchOne(ctx, repository.Query{
		Query:  "key_hash = ?",
		Values: []interface{}{keyHash},
	})
	if err != nil || !apiKey.IsActive {
		return dto.VerifyAPIKeyResponse{}, err
	}

	return dto.VerifyAPIKeyResponse{
		IsValid:            true,
		ServiceName:        apiKey.Name,
		ThresholdRateLimit: apiKey.MaxPerMinute,
		APIKey:             apiKey.KeyHash,
	}, nil
}
