package main

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
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

func writeJson(w http.ResponseWriter, statusCode int, data any) (int, error) {
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func allowableSizeOfUserInputs(w interface{}) error {
	switch v := w.(type) {
	case AccountSignIn:
		if len(v.UserName) > 20 || len(v.Password) > 20 {
			return errors.New("length of user <= 20 or password <= 20")
		}
	case AccountSignUp:
		if len(v.UserName) > 20 || len(v.Password) > 20 {
			return errors.New("length of user <= 20 or password <= 20")
		}
		if len(v.Name) > 50 || len(v.Email) > 30 {
			return errors.New("length of name <= 50 or email <= 30")
		}
	default:
		return errors.New("unsupported type")
	}
	return nil
}
