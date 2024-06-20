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
	"file-manager/src/handlers"
	"file-manager/src/serverConfigs"

	"github.com/labstack/echo/v4"
)



func main() {
	server := echo.New()

	serverConfigs.LoadEnvVariables(server)
	
	serverConfigs.SetupS3PresignClient(server)
	serverConfigs.ConnectToMongo()

	server.GET("/", handlers.HelloWorldHandler)

	server.POST("/presign", handlers.CreateUploadUrl)
	server.GET("/presign", handlers.GetDownloadUrl)
	server.GET("/folder", handlers.GetFolderContents)

	serverConfigs.StartListner(server)
}