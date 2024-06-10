package handlers

import (
	"auth/src/models"
	"auth/src/store/types"
	"auth/src/utils"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Login(c echo.Context) error {
	jsonBody := types.LoginInput{}
	c.Bind(&jsonBody)
	isValid, inputErr := jsonBody.Validate()

	if !isValid {

		println(inputErr.Error())
		responseData := types.ResponseBody{
			Message: "Error parsing json, please check type of each parameter",
			Data: inputErr.Error(),
		}
		return c.JSON(http.StatusBadRequest, responseData)
	}

	inputUser := models.UserModel{
		Email: jsonBody.Email,
		Phone: jsonBody.Phone,
		Password: jsonBody.Password,
	}

	userFetchResult, findError := inputUser.GetUser()

	if findError != nil {
		println(findError.Error())
		responseData := types.ResponseBody{
			Message: "Error finding the user",
			Data: findError.Error(),
		}
		return c.JSON(http.StatusBadRequest, responseData)
	}

	user := userFetchResult.Data.(models.UserModel)

	passwordValidationError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(jsonBody.Password))

	if passwordValidationError != nil {
		println(passwordValidationError.Error())
		responseData := types.ResponseBody{
			Message: "Wrong password!",
			Data: passwordValidationError.Error(),
		}
		return c.JSON(http.StatusBadRequest, responseData)
	}

	userDetails := jwt.MapClaims{
		"userId": user.Id,
		"orgId": user.Org,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	signedJwt, err := utils.SignJwtToken(&userDetails)

	if err != nil {
		println(err.Error())
		responseData := types.ResponseBody{
			Message: "Error signing token",
			Data: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, responseData)
	}

	responseData := types.ResponseBody{
		Message: "User loggedin successfully",
		Data: utils.SignedJwtToken{
			Token: signedJwt,
		},
	}
	return c.JSON(http.StatusOK, responseData)
}
