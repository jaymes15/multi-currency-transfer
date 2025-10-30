package bootstrap

import (
	"lemfi/simplebank/config"
	"lemfi/simplebank/db"
	"lemfi/simplebank/pkg/token"

	"github.com/joho/godotenv"
)

func Serve() {
	godotenv.Load()
	config.Set()
	db.Connect()
	token.SetTokenMaker()
	PostgresDB := db.GetPostgresDBConnection()

	// Start gRPC server in a goroutine so it runs in the background

	go GrpcGatewayServe()
	GrpcServe()

	//routing.RouteBuilder()

	defer PostgresDB.Close()
}
