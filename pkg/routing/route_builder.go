package routing

import (
	"os"

	"lemfi/simplebank/config"
)

func RouteBuilder() {
	Init()
	route := getRouter()
	registeredRoutes := registerRoutes(route)
	err := serve(registeredRoutes)
	config.Logger.Error(err.Error())
	os.Exit(1)
}
