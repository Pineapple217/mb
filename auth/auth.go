package auth

import (
	"crypto/rand"
	"math/big"
)

var SecretPassword string = generateRandomPassword(32)

var SecretCookie string

func generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	password := make([]byte, length)

	for i := range password {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "password"
		}
		password[i] = charset[randomIndex.Int64()]
	}

	return string(password)
}
