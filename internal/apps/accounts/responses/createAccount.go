package accounts

import (
	"time"

	"github.com/shopspring/decimal"
)

type CreateAccountResponse struct {
	ID        int64           `json:"id"`
	Owner     string          `json:"owner"`
	Balance   decimal.Decimal `json:"balance"`
	Currency  string          `json:"currency"`
	CreatedAt time.Time       `json:"created_at"`
}
