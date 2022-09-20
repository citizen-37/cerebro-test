package jsonrpc

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/bitwurx/jrpc2"

	"cerebro-test/internal/repositories"
	"cerebro-test/internal/services"
)

type Server struct {
	userRepo UserRepository
	service  TransactionService
}

func NewServer(userRepo UserRepository, service TransactionService) *Server {
	return &Server{
		userRepo: userRepo,
		service:  service,
	}
}

func (s *Server) GetBalance(ctx context.Context, params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	var request *GetBalanceRequest
	err := json.Unmarshal(params, request)
	if err != nil {
		return nil, newError(jrpc2.ParseErrorCode, err.Error())
	}

	balance, err := s.userRepo.GetBalance(ctx, request.PlayerName)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, newError(jrpc2.InvalidParamsCode, err.Error())
		}

		// FIXME: log
		return nil, newError(jrpc2.InternalErrorCode, "internal error")
	}

	return &GetBalanceResponse{Balance: balance}, nil
}

func (s *Server) Transaction(ctx context.Context, params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	var request WithdrawAndDepositRequest
	err := json.Unmarshal(params, &request)
	if err != nil {
		return nil, newError(jrpc2.ParseErrorCode, err.Error())
	}

	response, err := s.service.Create(ctx, &services.Transaction{
		CallerID:             request.CallerID,
		PlayerName:           request.PlayerName,
		Withdraw:             request.Withdraw,
		Deposit:              request.Deposit,
		Currency:             request.Currency,
		TransactionRef:       request.TransactionRef,
		GameRoundRef:         request.GameRoundRef,
		GameID:               request.GameID,
		Source:               request.Source,
		Reason:               string(request.Reason),
		SessionID:            request.SessionID,
		SessionAlternativeID: request.SessionAlternativeID,
		SpinDetails: &services.SpinDetails{
			BetType: request.SpinDetails.BetType,
			WinType: request.SpinDetails.BetType,
		},
		BonusID:          request.BonusID,
		ChargeFreeRounds: request.ChargeFreeRounds,
	})
	if err != nil {
		return nil, newError(jrpc2.InternalErrorCode, "internal error")
	}

	return &WithdrawAndDepositResponse{
		NewBalance:    response.NewBalance,
		TransactionID: response.TransactionID,
	}, nil
}

func (s *Server) Rollback(ctx context.Context, params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	var request RollbackTransactionRequest
	err := json.Unmarshal(params, &request)
	if err != nil {
		return nil, newError(jrpc2.ParseErrorCode, err.Error())
	}

	err = s.service.Rollback(ctx, request.PlayerName, request.TransactionRef)
	if err != nil {
		// FIXME: log
		return nil, newError(jrpc2.InternalErrorCode, "internal error")
	}

	return nil, nil
}

func newError(code jrpc2.ErrorCode, message string) *jrpc2.ErrorObject {
	return &jrpc2.ErrorObject{
		Code:    code,
		Message: jrpc2.ErrorMsg(message),
	}
}
