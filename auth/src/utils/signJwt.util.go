package utils

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func SignJwtToken(data *jwt.MapClaims) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
