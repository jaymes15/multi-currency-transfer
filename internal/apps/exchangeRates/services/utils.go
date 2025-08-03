package exchangeRates

import (
	"lemfi/simplebank/config"
	"time"

	db "lemfi/simplebank/db/sqlc"
)

func (exchangeRateService *ExchangeRateService) IsExchangeRateExpired(exchangeRate db.ExchangeRate) bool {
	config.Logger.Info("Checking if exchange rate is expired", "exchange_rate", exchangeRate)

	// Calculate when this exchange rate expires
	expiredTime := exchangeRateService.GetExchangeRateExpiredTime(exchangeRate.UpdatedAt.Time)

	// Check if current time is after the expiration time
	isExpired := time.Now().After(expiredTime)

	config.Logger.Info("Exchange rate expiration check",
		"created_at", exchangeRate.CreatedAt.Time,
		"updated_at", exchangeRate.UpdatedAt.Time,
		"expired_time", expiredTime,
		"current_time", time.Now(),
		"is_expired", isExpired,
	)

	return isExpired
}

func (exchangeRateService *ExchangeRateService) GetExchangeRateExpiredTime(updatedAt time.Time) time.Time {
	cfg := config.Get()
	expiredTime := updatedAt.Add(time.Duration(cfg.ExchangeRate.ExpiredTimeInMinutes) * time.Minute)

	config.Logger.Info("Calculating expired time",
		"updated_at", updatedAt,
		"expired_time_in_minutes", cfg.ExchangeRate.ExpiredTimeInMinutes,
		"expired_time", expiredTime,
	)

	return expiredTime
}
