package exchangeRates

import (
	"context"

	dbConnection "lemfi/simplebank/db"
	db "lemfi/simplebank/db/sqlc"
)

type ExchangeRateRepository struct {
	context context.Context
	queries db.Store
}

func NewExchangeRateRepository() *ExchangeRateRepository {
	return &ExchangeRateRepository{
		context: context.Background(),
		queries: db.NewStore(dbConnection.GetPostgresDBConnection()),
	}
}
