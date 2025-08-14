package bootstrap

import (
	"lemfi/simplebank/config"
	"lemfi/simplebank/db"
	"lemfi/simplebank/pkg/routing"
	"lemfi/simplebank/pkg/token"

	"github.com/joho/godotenv"
)

func Serve() {
	godotenv.Load()
	config.Set()
	db.Connect()
	token.SetTokenMaker()
	PostgresDB := db.GetPostgresDBConnection()

	routing.RouteBuilder()

	defer PostgresDB.Close()

}
