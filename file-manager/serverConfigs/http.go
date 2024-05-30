package serverConfig

import (
	"os"

	"github.com/labstack/echo/v4"
)

func StartListner(server *echo.Echo) {

	port := ":" + os.Getenv("PORT")

	if port == ":" {
		server.Logger.Fatal("Error starting server:", "Port is missing in env")
	}

	err := server.Start(port)

	if err != nil {
		server.Logger.Fatal("Error starting server:", err)
	}
}
