package jsonrpc

const (
	ReasonGamePlay      Reason = "GAME_PLAY"
	ReasonGamePlayFinal Reason = "GAME_PLAY_FINAL"
)

type (
	Reason string

	GetBalanceRequest struct {
		CallerID   int64  `json:"callerId"`
		PlayerName string `json:"playerName"`
		Currency   string `json:"currency"`

		GameID               string `json:"gameId,omitempty"`
		SessionID            string `json:"sessionId,omitempty"`
		SessionAlternativeID string `json:"sessionAlternativeId,omitempty"`
		BonusID              string `json:"bonusId,omitempty"`
	}

	GetBalanceResponse struct {
		Balance        int64 `json:"balance"`
		FreeRoundsLeft int   `json:"freeroundsLeft,omitempty"`
	}

	WithdrawAndDepositRequest struct {
		CallerID             int64        `json:"callerId,omitempty"`
		PlayerName           string       `json:"playerName,omitempty"`
		Withdraw             int64        `json:"withdraw,omitempty"`
		Deposit              int64        `json:"deposit,omitempty"`
		Currency             string       `json:"currency,omitempty"`
		TransactionRef       string       `json:"transactionRef,omitempty"`
		GameRoundRef         string       `json:"gameRoundRef,omitempty"`
		GameID               string       `json:"gameId,omitempty"`
		Source               string       `json:"source,omitempty"`
		Reason               Reason       `json:"reason,omitempty"`
		SessionID            string       `json:"sessionId,omitempty"`
		SessionAlternativeID string       `json:"sessionAlternativeId,omitempty"`
		SpinDetails          *SpinDetails `json:"spinDetails,omitempty"`
		BonusID              string       `json:"bonusID,omitempty"`
		ChargeFreeRounds     int64        `json:"chargeFreerounds,omitempty"`
	}

	SpinDetails struct {
		BetType string `json:"betType,omitempty"`
		WinType string `json:"winType,omitempty"`
	}

	WithdrawAndDepositResponse struct {
		NewBalance     int64  `json:"newBalance,omitempty"`
		TransactionID  string `json:"transactionId,omitempty"`
		FreeRoundsLeft int64  `json:"freeroundsLeft,omitempty"`
	}

	RollbackTransactionRequest struct {
		CallerID             int64  `json:"callerId,omitempty"`
		PlayerName           string `json:"playerName,omitempty"`
		TransactionRef       string `json:"transactionRef,omitempty"`
		GameID               string `json:"gameId,omitempty"`
		SessionID            string `json:"sessionId,omitempty"`
		SessionAlternativeID string `json:"sessionAlternativeId,omitempty"`
		RoundID              string `json:"roundId,omitempty"`
	}

	RollbackTransactionResponse struct{}
)
