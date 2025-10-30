package config

import (
	"flag"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

func Set() Config {

	exchangeRateExpiredTimeInMinutes, err := strconv.Atoi(os.Getenv("EXCHANGE_RATE_EXPIRED_TIME_IN_MINUTES"))
	if err != nil {
		Logger.Error("Failed to convert EXCHANGE_RATE_EXPIRED_TIME_IN_MINUTES to int", "error", err.Error())
		panic(err)
	}

	multiCurrencyFee, err := decimal.NewFromString(os.Getenv("MULTI_CURRENCY_FEE"))
	if err != nil {
		// Default to 1.99 if environment variable is not set or invalid
		multiCurrencyFee = decimal.NewFromFloat(1.99)
		Logger.Info("Using default multi-currency fee", "fee", multiCurrencyFee.String())
	}

	// Create a string variable to hold the fee flag value
	var feeFlag string

	// Set configurations using environment variables or flags
	flag.IntVar(&configurations.Port, "port", 4000, "API server port")
	flag.StringVar(&configurations.Env, "env", os.Getenv("ENVIROMENT"), "Environment (development|staging|production)")
	flag.StringVar(&configurations.Db.Dsn, "db-dsn", os.Getenv("POSTGRES_DB_CONNECTION_STRING"), "PostgreSQL DSN")
	flag.IntVar(&configurations.Db.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&configurations.Db.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&configurations.Db.MaxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")
	flag.IntVar(&configurations.ExchangeRate.ExpiredTimeInMinutes, "exchange-rate-expired-time-in-minutes", exchangeRateExpiredTimeInMinutes, "Exchange rate expired time in minutes")
	flag.StringVar(&feeFlag, "multi-currency-fee", multiCurrencyFee.String(), "Multi currency fee")
	flag.StringVar(&configurations.TokenSymmetricKey, "token-symmetric-key", os.Getenv("TOKEN_SYMMETRIC_KEY"), "Token symmetric key")
	flag.DurationVar(&configurations.AccessTokenDuration, "access-token-duration", 15*time.Minute, "Access token duration")
	flag.DurationVar(&configurations.RefreshTokenDuration, "refresh-token-duration", 7*24*time.Hour, "Refresh token duration")
	flag.StringVar(&configurations.GRPCServerAddress, "grpc-server-address", os.Getenv("GRPC_SERVER_ADDRESS"), "gRPC server address")

	// Parse the flags
	flag.Parse()

	// Convert the fee flag string to decimal
	if feeFlag != "" {
		fee, err := decimal.NewFromString(feeFlag)
		if err != nil {
			Logger.Error("Failed to convert multi-currency-fee flag to decimal", "error", err.Error())
			panic(err)
		}
		configurations.MultiCurrency.Fee = fee
	} else {
		configurations.MultiCurrency.Fee = multiCurrencyFee
	}

	// Set CORS Trusted Origins
	configurations.Cors.TrustedOrigins = strings.Fields(os.Getenv("TRUSTED_ORIGINS"))

	// Set default gRPC server address if not provided
	if configurations.GRPCServerAddress == "" {
		configurations.GRPCServerAddress = ":9090"
	}

	return configurations
}
