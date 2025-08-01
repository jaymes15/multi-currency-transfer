package bootstrap

import (
	"lemfi/simplebank/config"
	"lemfi/simplebank/db"
	"lemfi/simplebank/pkg/routing"

	"github.com/joho/godotenv"
)

func Serve() {
	godotenv.Load()
	config.Set()
	db.Connect()
	PostgresDB := db.GetPostgresDBConnection()

	routing.RouteBuilder()

	defer PostgresDB.Close()

}
