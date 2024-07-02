package types

import (
	"errors"
	"strings"
)

type ValidateTokenInput struct {
	Token string `headers:"authorization"`
}

func (tokenInput *ValidateTokenInput) Validate() (string, error) {
	if tokenInput.Token == "" {
		return "", errors.New("authorization failed, token is missing")
	}

	bearerToken := tokenInput.Token
	tokenSplitArray := strings.Split(bearerToken, " ")

	token := tokenSplitArray[1]

	if token == "" {
		return "", errors.New("authorization failed, invalid token")
	}

	return token, nil
}