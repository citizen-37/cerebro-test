package services

import (
	"context"

	"github.com/pkg/errors"

	"cerebro-test/internal/models"
)

type TransactionService struct {
	transactionRepo TransactionRepository
	userRepo        UserRepository
}

func NewTransactionService(userRepo UserRepository, transactionRepo TransactionRepository) *TransactionService {
	return &TransactionService{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *TransactionService) Create(ctx context.Context, data *Transaction) (*TransactionResponse, error) {
	amount := data.Deposit
	if data.Withdraw != 0 {
		amount = -data.Withdraw
	}

	instance := &models.Transaction{
		UserName:   data.PlayerName,
		ExternalID: data.TransactionRef,
		Amount:     amount,
		Payload: &models.TransactionPayload{
			Currency:             data.Currency,
			GameRoundRef:         data.GameRoundRef,
			GameID:               data.GameID,
			Source:               data.Source,
			Reason:               data.Reason,
			SessionID:            data.SessionID,
			SessionAlternativeID: data.SessionAlternativeID,
			SpinDetails: &models.SpinDetails{
				BetType: data.SpinDetails.BetType,
				WinType: data.SpinDetails.WinType,
			},
			BonusID:          data.BonusID,
			ChargeFreeRounds: data.ChargeFreeRounds,
		},
	}

	err := s.transactionRepo.Create(ctx, instance)
	if err != nil {
		return nil, errors.Wrap(err, "cant create transaction")
	}

	balance, err := s.userRepo.GetBalance(ctx, data.PlayerName)
	if err != nil {
		return nil, errors.Wrap(err, "cant get user balance")
	}

	return &TransactionResponse{
		NewBalance:    balance,
		TransactionID: instance.ID,
	}, nil
}

func (s *TransactionService) Rollback(ctx context.Context, userName, externalID string) error {
	return s.transactionRepo.Rollback(ctx, userName, externalID)
}
