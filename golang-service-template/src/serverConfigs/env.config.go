package serverConfigs

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func LoadEnvVariables(server *echo.Echo) {
	err := godotenv.Load()

	if err != nil {
		server.Logger.Fatal("Error loading .env file:", err)
	}
}
