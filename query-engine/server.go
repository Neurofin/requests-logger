/**
* Author: Sree
* Just some Go code magic happening here!
* Warning: Extreme levels of awesomeness detected.
* Remember, great code comes with great responsibility.
* If you find any bugs, they’re probably just features in disguise!
* Keep coding, stay curious, and have fun!
**/

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