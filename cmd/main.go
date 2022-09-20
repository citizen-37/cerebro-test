package main

import (
	"fmt"

	"github.com/bitwurx/jrpc2"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/kelseyhightower/envconfig"

	"cerebro-test/internal/config"
	"cerebro-test/internal/jsonrpc"
	"cerebro-test/internal/repositories"
	"cerebro-test/internal/services"
)

func main() {
	cfg := &config.Config{}
	if err := envconfig.Process("CEREBRO_TEST", cfg); err != nil {
		panic(fmt.Errorf("error reading config: %w", err))
	}

	db, err := sqlx.Connect("pgx", cfg.PG.DSN)
	if err != nil {
		panic(errors.Wrap(err, "cant PG connect"))
	}

	userRepo := repositories.NewUserRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(userRepo, transactionRepo)

	app := jsonrpc.NewServer(userRepo, transactionService)

	server := jrpc2.NewServer(":9090", "/", nil)
	server.RegisterWithContext("getBalance", jrpc2.MethodWithContext{Method: app.GetBalance})
	server.RegisterWithContext("withdrawAndDeposit", jrpc2.MethodWithContext{Method: app.Transaction})
	server.RegisterWithContext("rollbackTransaction", jrpc2.MethodWithContext{Method: app.Rollback})

	server.Start()
}
