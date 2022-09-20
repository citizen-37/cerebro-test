package services

import (
	"context"

	"cerebro-test/internal/models"
)

type (
	UserRepository interface {
		GetBalance(ctx context.Context, name string) (int64, error)
	}

	TransactionRepository interface {
		Create(ctx context.Context, instance *models.Transaction) error
		Rollback(ctx context.Context, userName, externalID string) error
	}
)
