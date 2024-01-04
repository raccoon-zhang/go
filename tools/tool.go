package tools

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func PasswordEncrypt(password string) (string, error) {
	password_encryption, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(password_encryption), err
}

func PasswordDecrypt(hashPassWord, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassWord), []byte(password))
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
