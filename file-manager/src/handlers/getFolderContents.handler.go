package handlers

import (
	"file-manager/src/serverConfigs"
	"net/http"

	"github.com/labstack/echo/v4"
)

type GetFolderContentsInput struct {
	FolderPath string `json:"folderPath" query:"folderPath"`
	BucketName string `json:"bucketName" query:"bucketName"`
} 

func GetFolderContents(c echo.Context) error {

	input := GetFolderContentsInput{}

	inputErr := c.Bind(&input)
	if inputErr != nil {

		println(inputErr.Error())
		responseData := ResponseBody{
			Message: "Error in input params, please verify them again",
			Data: inputErr.Error(),
		}
		return c.JSON(http.StatusBadRequest, responseData)

	}

	results, err := serverConfigs.ListS3Objects(input.BucketName, input.FolderPath)
	if err != nil {
		println(err.Error())
		responseData := ResponseBody{
			Message: "Couldn't get a folder contents",
			Data: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, responseData)
	}

	responseData := ResponseBody{
		Message: "Fetched folder contents successfully",
		Data: results,
	}
	return c.JSON(http.StatusOK, responseData)
}
