package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// bcrypt.DefaultCost = 10 (you can use 12 or more for stronger hash)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

