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
	"logger/src/handlers"
	"logger/src/serverConfigs"

	"github.com/labstack/echo/v4"
)

func main() {
	server := echo.New()

	serverConfigs.LoadEnvVariables(server)

	server.GET("/", handlers.HelloWorldHandler)

	serverConfigs.StartListner(server)
}
