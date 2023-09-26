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

func GenerateRandomString() (string, error) {
	const n int = 5
	const letters string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

func main() {
	s := "1234"
	a := "1234"

	salt1, err := GenerateRandomString()
	if err != nil {
		panic(err)
	}

	salt2, err := GenerateRandomString()
	if err != nil {
		panic(err)
	}

	fmt.Println(genHash(s + salt1))
	fmt.Println(genHash(a + salt2))
}
