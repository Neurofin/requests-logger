package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func EncrytPassword(password string) (string, error) {
	encryptedPassword, encryptionError := bcrypt.GenerateFromPassword([]byte(password), 10)

	if encryptionError != nil {
		println(encryptionError.Error())
		return "", encryptionError
	}

	stringifiedEncryptedPassword := string(encryptedPassword[:])
	return stringifiedEncryptedPassword, nil
}
