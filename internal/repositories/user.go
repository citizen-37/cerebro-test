package repositories

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"

	"cerebro-test/internal/models"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetBalance(ctx context.Context, name string) (int64, error) {
	result := &models.User{}

	err := r.db.GetContext(ctx, result, "select balance from users where name = ?", name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrUserNotFound
		}
	}

	return result.Balance, nil
}
