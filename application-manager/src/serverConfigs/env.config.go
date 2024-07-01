package serverConfigs

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

var ApplicationDocumentsBucket string

func LoadEnvVariables(server *echo.Echo) {
	err := godotenv.Load()

	if err != nil {
		server.Logger.Fatal("Error loading .env file:", err)
	}

	ApplicationDocumentsBucket = os.Getenv("DOCUMENTS_BUCKET")
}
