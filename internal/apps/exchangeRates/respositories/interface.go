package exchangeRates

import (
	"context"

	db "lemfi/simplebank/db/sqlc"
	requests "lemfi/simplebank/internal/apps/exchangeRates/requests"
)

type ExchangeRateRepositoryInterface interface {
	ListExchangeRates(ctx context.Context) ([]db.ExchangeRate, error)
	GetExchangeRate(ctx context.Context, payload requests.GetExchangeRateRequest) (db.ExchangeRate, error)
}
