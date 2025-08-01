package db

import (
	"context"
	"database/sql"
	"lemfi/simplebank/config"
	"time"

	_ "github.com/lib/pq"
)

func ConnectToPostgresDb() (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config
	// struct.
	config.Logger.Info("Connecting to Database")
	configs := config.Get()
	db, err := sql.Open("postgres", configs.Db.Dsn)
	if err != nil {
		return nil, err
	}

	// Set the maximum number of open (in-use + idle) connections in the pool. Note that
	// passing a value less than or equal to 0 will mean there is no limit.
	db.SetMaxOpenConns(configs.Db.MaxOpenConns)

	// Set the maximum number of idle connections in the pool. Again, passing a value
	// less than or equal to 0 will mean there is no limit.
	db.SetMaxIdleConns(configs.Db.MaxIdleConns)

	// Set the maximum idle timeout for connections in the pool. Passing a duration less
	// than or equal to 0 will mean that connections are not closed due to their idle time.
	db.SetConnMaxIdleTime(configs.Db.MaxIdleTime)

	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use PingContext() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error. If we get this error, or any other, we close the connection pool and
	// return the error.
	err = db.PingContext(ctx)
	if err != nil {
		config.Logger.Error(err.Error())
		db.Close()
		return nil, err
	}

	config.Logger.Info("Sucessfully  connected to Database")

	// Return the sql.DB connection pool.
	return db, nil

}
