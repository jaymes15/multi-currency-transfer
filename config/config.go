package config

import (
	"time"

	"github.com/shopspring/decimal"
)

type Config struct {
	Port int
	Env  string
	Db   struct {
		Dsn          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  time.Duration
	}
	Cors struct {
		TrustedOrigins []string
	}
	ExchangeRate struct {
		ExpiredTimeInMinutes int
	}
	MultiCurrency struct {
		Fee decimal.Decimal
	}
	TokenSymmetricKey    string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	GRPCServerAddress    string
}
