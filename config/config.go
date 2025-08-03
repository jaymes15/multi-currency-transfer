package config

import "time"

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
}
