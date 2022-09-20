package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"cerebro-test/internal/models"
)

type TransactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (r *TransactionRepository) Create(ctx context.Context, instance *models.Transaction) error {
	var err error
	var userID string

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "cant start transaction")
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback() // FIXME: log
		}
	}()

	query := `select id from users where name = $1 for update`

	err = tx.QueryRowContext(ctx, query, instance.UserName).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}

		return errors.Wrap(err, "cant select user for update")
	}

	instance.ID = uuid.New().String()

	_, err = tx.ExecContext(
		ctx,
		"insert into transactions (id, user_id, external_id, amount, created_at, payload) values ($1, $2, $3, $4, $5, $6)",
		instance.ID,
		userID,
		instance.ExternalID,
		instance.Amount,
		time.Now().UTC(),
		instance.Payload,
	)
	if err != nil {
		return errors.Wrap(err, "cant insert transaction")
	}

	_, err = tx.ExecContext(ctx, "update users set balance = balance + $2 where id = $1", userID, instance.Amount)
	if err != nil {
		return errors.Wrap(err, "cant update balance")
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "cant commit transaction")
	}

	return nil
}

func (r *TransactionRepository) Rollback(ctx context.Context, userName, externalID string) error {
	var err error
	var userID string

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "cant start transaction")
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback() // FIXME: log
		}
	}()

	query := `select id from users where name = $1 for update`

	err = tx.QueryRowContext(ctx, query, userName).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}

		return errors.Wrap(err, "cant select user for update")
	}

	instance := &models.Transaction{}

	err = tx.QueryRowxContext(
		ctx,
		"select id, amount from transactions where external_id = $1",
		externalID,
	).StructScan(instance)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return errors.Wrap(err, "cant get transaction")
		}
	}

	if instance.ID != "" {
		_, err = tx.ExecContext(
			ctx,
			"update transactions set rolled_back_at = $2 where id = $1",
			instance.ID,
			time.Now().UTC(),
		)
		if err != nil {
			return errors.Wrap(err, "cant update transaction")
		}

		_, err = tx.ExecContext(
			ctx,
			"update users set balance = balance - $2 where id = $1",
			userID,
			instance.Amount,
		)
		if err != nil {
			return errors.Wrap(err, "cant update balance")
		}
	} else {
		_, err = tx.ExecContext(
			ctx,
			"insert into transactions (id, external_id, rolled_back_at) values ($1, $2, $3)",
			uuid.New(),
			externalID,
			time.Now().UTC(),
		)
		if err != nil {
			return errors.Wrap(err, "cant insert transaction")
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "cant commit transaction")
	}

	return nil
}
