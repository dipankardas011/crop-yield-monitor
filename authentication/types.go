package main

import (
	"fmt"

	jwt "github.com/golang-jwt/jwt/v5"
)

type AccountSignUp struct {
	Name     string `json:"name"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"emailid"`
}

type AccountSignInRes struct {
	Uuid string `json:"uuid"`
	// AccessToken string    `json:"token"`
	// ExpTime     time.Time `json:"expiration"`
}

type AccountSignIn struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Stdout  string `json:"stdout"`
	Error   string `json:"error"`
	Account any
}

func (this AccountSignUp) String() string {
	return fmt.Sprintf("{name: %s, username: %s, password: %s, email: %s}\n", this.Name, this.UserName, this.Password, this.Email)
}

func (this AccountSignIn) String() string {
	return fmt.Sprintf("{username: %s, password: %s}\n", this.UserName, this.Password)
}

// Create a struct that will be encoded to a JWT.
// We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
