package responses

import (
	"lemfi/simplebank/config"
	db "lemfi/simplebank/db/sqlc"
	"time"

	"github.com/shopspring/decimal"
)

type ExchangeRateResponse struct {
	ID           int64           `json:"id"`
	FromCurrency string          `json:"from_currency"`
	ToCurrency   string          `json:"to_currency"`
	Rate         decimal.Decimal `json:"rate"`
	CreatedAt    time.Time       `json:"created_at"`
	ExpiredAt    time.Time       `json:"expired_at"`
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

// NewExchangeRateResponse creates a new ExchangeRateResponse with calculated expired time
func NewExchangeRateResponse(dbExchangeRate db.ExchangeRate) ExchangeRateResponse {
	cfg := config.Get()
	expiredAt := dbExchangeRate.CreatedAt.Time.Add(time.Duration(cfg.ExchangeRate.ExpiredTimeInMinutes) * time.Minute)

	return ExchangeRateResponse{
		ID:           dbExchangeRate.ID,
		FromCurrency: dbExchangeRate.FromCurrency,
		ToCurrency:   dbExchangeRate.ToCurrency,
		Rate:         dbExchangeRate.Rate,
		CreatedAt:    dbExchangeRate.CreatedAt.Time,
		ExpiredAt:    expiredAt,
	}
}
