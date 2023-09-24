package main

import (
	"math/rand"
	"strings"
	"time"
)

func generateRandomToken(length int) string {
	var password strings.Builder
	var (
		lowerCharSet = "abcdedfghijklmnopqrst"
		upperCharSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		numberSet    = "0123456789"
		special      = "@$#{}"
		allCharSet   = lowerCharSet + upperCharSet + numberSet + special
	)
	rand.Seed(time.Now().Unix())

	for i := 0; i < length; i++ {
		random := rand.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[random]))
	}

	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})

	return string(inRune)
}
