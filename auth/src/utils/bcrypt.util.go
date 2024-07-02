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

func VerifyPassword(inputPassword string, encryptedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(inputPassword))

	if err != nil {
		return err
	}

	return nil
}
