package config

import (
	"flag"
	"os"
	"strconv"
	"strings"
	"time"
)

func Set() Config {

	exchangeRateExpiredTimeInMinutes, err := strconv.Atoi(os.Getenv("EXCHANGE_RATE_EXPIRED_TIME_IN_MINUTES"))
	if err != nil {
		Logger.Error("Failed to convert EXCHANGE_RATE_EXPIRED_TIME_IN_MINUTES to int", "error", err.Error())
		panic(err)
	}

	// Set configurations using environment variables or flags
	flag.IntVar(&configurations.Port, "port", 4000, "API server port")
	flag.StringVar(&configurations.Env, "env", os.Getenv("ENVIROMENT"), "Environment (development|staging|production)")
	flag.StringVar(&configurations.Db.Dsn, "db-dsn", os.Getenv("POSTGRES_DB_CONNECTION_STRING"), "PostgreSQL DSN")
	flag.IntVar(&configurations.Db.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&configurations.Db.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&configurations.Db.MaxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")
	flag.IntVar(&configurations.ExchangeRate.ExpiredTimeInMinutes, "exchange-rate-expired-time-in-minutes", exchangeRateExpiredTimeInMinutes, "Exchange rate expired time in minutes")

	// Parse the flags
	flag.Parse()

	// Set CORS Trusted Origins
	configurations.Cors.TrustedOrigins = strings.Fields(os.Getenv("TRUSTED_ORIGINS"))

	return configurations
}
