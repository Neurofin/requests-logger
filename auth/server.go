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
	
	server.POST("admin/org", handlers.CreateOrg)
	server.GET("admin/org", handlers.GetOrg)

	server.POST("admin/user/signup", handlers.AdminSignup)
	server.POST("/user/login", handlers.Login)

	server.Use(serverMiddleware.ValidateToken)
	server.POST("/user/signup", handlers.Signup)
	server.GET("/user/validate", handlers.ValidateToken)

	serverConfigs.StartListner(server)
}