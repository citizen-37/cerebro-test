package models

import (
	"database/sql/driver"
	"encoding/json"
)

type (
	User struct {
		Balance int64 `db:"balance"`
	}

	Transaction struct {
		ID         string              `db:"id"`
		UserName   string              `db:"user_name"`
		ExternalID string              `db:"external_id"`
		Amount     int64               `db:"amount"`
		Payload    *TransactionPayload `db:"payload"`
	}

	TransactionPayload struct {
		Currency             string       `json:"currency,omitempty"`
		GameRoundRef         string       `json:"gameRoundRef,omitempty"`
		GameID               string       `json:"gameId,omitempty"`
		Source               string       `json:"source,omitempty"`
		Reason               string       `json:"reason,omitempty"`
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
)

func (p *TransactionPayload) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *TransactionPayload) Scan(value interface{}) error {
	data := value.([]byte)
	return json.Unmarshal(data, p)
}
