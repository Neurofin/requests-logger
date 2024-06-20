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

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)



func main() {
	server := echo.New()

	serverConfigs.LoadEnvVariables(server)

	serverConfigs.ConnectToMongo()

	server.Use(middleware.Logger())

	server.GET("/", handlers.HelloWorldHandler)
	
	// server.POST("admin/org", handlers.CreateOrg)
	// server.GET("admin/org", handlers.GetOrg)

	server.POST("/admin/user/signup", handlers.AdminSignup)
	server.POST("/user/login", handlers.Login)

	server.POST("/user/signup", handlers.Signup, serverMiddleware.ValidateToken)
	server.GET("/user/validate", handlers.ValidateToken, serverMiddleware.ValidateToken)

	serverConfigs.StartListner(server)
}