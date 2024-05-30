package main

import (
	handler "file-manager/handlers"
	serverConfig "file-manager/serverConfigs"

	"github.com/labstack/echo/v4"
)



func main() {
	server := echo.New()

	serverConfig.LoadEnvVariables(server)
	
	server.GET("/", handler.HelloWorldHandler)

	serverConfig.StartListner(server)
}