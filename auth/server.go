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
	"auth/src/handlers"
	serverMiddleware "auth/src/middleware"
	"auth/src/serverConfigs"
	logger "github.com/Neurofin/requests-logger/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	server := echo.New()

	serverConfigs.LoadEnvVariables(server)

	serverConfigs.ConnectToMongo()

	server.Use(middleware.Logger())
	server.Use(middleware.CORS())
	server.Use(logger.LoggingMiddleware)

	server.GET("/auth", handlers.HelloWorldHandler)

	server.POST("/auth/admin/user/signup", handlers.AdminSignup)
	server.POST("/auth/user/login", handlers.Login)

	server.POST("/auth/user/signup", handlers.Signup, serverMiddleware.ValidateToken)
	server.GET("/auth/user/validate", handlers.ValidateToken, serverMiddleware.ValidateToken)

	serverConfigs.StartListner(server)
}
