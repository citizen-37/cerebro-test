package jsonrpc

import (
	"context"

	"cerebro-test/internal/services"
)

type (
	UserRepository interface {
		GetBalance(ctx context.Context, name string) (int64, error)
	}

	TransactionService interface {
		Create(ctx context.Context, data *services.Transaction) (*services.TransactionResponse, error)
		Rollback(ctx context.Context, userName, externalID string) error
	}
)
