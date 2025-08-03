package exchangeRates

import (
	"lemfi/simplebank/config"
	"time"

	db "lemfi/simplebank/db/sqlc"
)

func (exchangeRateService *ExchangeRateService) IsExchangeRateExpired(exchangeRate db.ExchangeRate) bool {
	config.Logger.Info("Checking if exchange rate is expired", "exchange_rate", exchangeRate)

	// Calculate when this exchange rate expires
	expiredTime := exchangeRateService.GetExchangeRateExpiredTime(exchangeRate.CreatedAt.Time)

	// Check if current time is after the expiration time
	return time.Now().After(expiredTime)
}

func (exchangeRateService *ExchangeRateService) GetExchangeRateExpiredTime(createdAt time.Time) time.Time {
	cfg := config.Get()
	return createdAt.Add(time.Duration(cfg.ExchangeRate.ExpiredTimeInMinutes) * time.Minute)
}
