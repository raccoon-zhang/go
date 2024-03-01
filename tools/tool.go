package tools

import (
	"fmt"
	"os"
	"path/filepath"

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

func SaveCacheFile(fileName string, cacheData []byte) (err error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return
	}
	cacheFilePath := filepath.Join(cacheDir, fileName)

	err = os.WriteFile(cacheFilePath, cacheData, 0644)
	return
}

func GetCacheFile(fileName string) (data []byte, err error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return
	}
	cacheFilePath := filepath.Join(cacheDir, fileName)
	data, err = os.ReadFile(cacheFilePath)
	return
}
