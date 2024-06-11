package main

import (
	"query-engine/src/handlers"
	serverMiddleware "query-engine/src/middleware"
	"query-engine/src/serverConfigs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)



func main() {
	server := echo.New()

	serverConfigs.LoadEnvVariables(server)

	serverConfigs.ConnectToMongo()

	server.Use(middleware.Logger())

	server.GET("/", handlers.HelloWorldHandler, serverMiddleware.ValidateToken)

	serverConfigs.StartListner(server)
}