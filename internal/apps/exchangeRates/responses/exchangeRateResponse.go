package responses

import (
	"time"

	"github.com/shopspring/decimal"
)

type ExchangeRateResponse struct {
	ID           int64           `json:"id"`
	FromCurrency string          `json:"from_currency"`
	ToCurrency   string          `json:"to_currency"`
	Rate         decimal.Decimal `json:"rate"`
	CreatedAt    time.Time       `json:"created_at"`
}

type ListExchangeRatesResponse struct {
	ExchangeRates []ExchangeRateResponse `json:"exchange_rates"`
	Total         int                    `json:"total"`
}

type GetExchangeRateResponse struct {
	ExchangeRate    ExchangeRateResponse `json:"exchange_rate"`
	AmountToSend    decimal.Decimal      `json:"amount_to_send"`
	AmountToReceive decimal.Decimal      `json:"amount_to_receive"`
	CanTransact     bool                 `json:"can_transact"`
	Message         string               `json:"message"`
}
