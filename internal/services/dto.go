package services

type (
	Transaction struct {
		CallerID             int64
		PlayerName           string
		Withdraw             int64
		Deposit              int64
		Currency             string
		TransactionRef       string
		GameRoundRef         string
		GameID               string
		Source               string
		Reason               string
		SessionID            string
		SessionAlternativeID string
		SpinDetails          *SpinDetails
		BonusID              string
		ChargeFreeRounds     int64
	}

	SpinDetails struct {
		BetType string
		WinType string
	}

	TransactionResponse struct {
		NewBalance    int64
		TransactionID string
	}
)
