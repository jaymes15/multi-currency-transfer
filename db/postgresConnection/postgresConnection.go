package db

import (
	"context"
	"lemfi/simplebank/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectToPostgresDb() (*pgxpool.Pool, error) {
	config.Logger.Info("Connecting to Database")
	configs := config.Get()

	// Create a context with a timeout for connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Parse the pgxpool config from the DSN
	pgxConfig, err := pgxpool.ParseConfig(configs.Db.Dsn)
	if err != nil {
		config.Logger.Error("Failed to parse database DSN: " + err.Error())
		return nil, err
	}

	// Apply pool settings
	pgxConfig.MaxConns = int32(configs.Db.MaxOpenConns) // MaxOpenConns in your config
	pgxConfig.MinConns = int32(configs.Db.MaxIdleConns) // MaxIdleConns in your config
	pgxConfig.MaxConnIdleTime = configs.Db.MaxIdleTime  // time.Duration

	// Connect to the database
	dbpool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		config.Logger.Error("Failed to connect to database: " + err.Error())
		return nil, err
	}

	// Ping to ensure connection is working
	if err := dbpool.Ping(ctx); err != nil {
		config.Logger.Error("Database ping failed: " + err.Error())
		dbpool.Close()
		return nil, err
	}

	config.Logger.Info("Successfully connected to Database")
	return dbpool, nil
}
