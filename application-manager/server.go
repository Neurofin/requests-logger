package main

import (
	"application-manager/src/handlers"
	serverMiddleware "application-manager/src/middleware"
	"application-manager/src/serverConfigs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	server := echo.New()

	serverConfigs.LoadEnvVariables(server)

	serverConfigs.ConnectToMongo()

	server.Use(middleware.Logger())

	server.GET("/", handlers.HelloWorldHandler, serverMiddleware.ValidateToken)
	server.GET("/application/view", handlers.ApplicationView)

	serverConfigs.StartListner(server)
}
