package unitofwork

import (
	"context"
	"database/sql"
	"worker-service/internal/repository"

	"gorm.io/gorm"
)

type uowStore struct {
	emailHistories repository.EmailHistoryRepository
}

type UnitOfWorkStore interface {
	EmailHistories() repository.EmailHistoryRepository
}

func (u uowStore) EmailHistories() repository.EmailHistoryRepository {
	return u.emailHistories
}

type unitOfWork struct {
	db *gorm.DB
}

type UnitOfWorkBlock func(UnitOfWorkStore) error

type UnitOfWork interface {
	Do(ctx context.Context, fn UnitOfWorkBlock, option sql.TxOptions) error
}

func NewUoW(db *gorm.DB) UnitOfWork {
	return &unitOfWork{db: db}
}

func (s *unitOfWork) Do(ctx context.Context, fn UnitOfWorkBlock, option sql.TxOptions) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		newStore := &uowStore{
			emailHistories: repository.NewEmailHistoryRepository(tx),
		}
		return fn(newStore)
	}, &option)
}
