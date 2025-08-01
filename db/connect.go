package db

import (
	"lemfi/simplebank/config"
	db "lemfi/simplebank/db/postgresConnection"
	"os"
)

func Connect() {
	postgresDB, err := db.ConnectToPostgresDb()
	if err != nil {
		config.Logger.Error(err.Error())
		os.Exit(1)
	}

	PostgresDB = postgresDB

}
