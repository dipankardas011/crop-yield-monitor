package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

func genHash(msg string) string {
	h := sha256.New()

	h.Write([]byte(msg))

	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

func generateRandomString() (string, error) {
	const lenOfSalt int = 6
	const letters string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, lenOfSalt)
	for i := 0; i < lenOfSalt; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

// GenerateHashForPasswordAndSalt returns hashed (password+salt), salt, and error
func GenerateHashForPasswordAndSalt(password string) (string, string, error) {

	salt, err := generateRandomString()
	if err != nil {
		return "", "", err
	}

	return genHash(password + salt), salt, nil
}
