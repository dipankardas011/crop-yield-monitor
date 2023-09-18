package main

import "fmt"

type AccountSignUp struct {
	Name     string `json:"name"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"emailid"`
}

type AccountSignIn struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (this AccountSignUp) String() string {
	return fmt.Sprintf("{name: %s, username: %s, password: %s, email: %s}\n", this.Name, this.UserName, this.Password, this.Email)
}

func (this AccountSignIn) String() string {
	return fmt.Sprintf("{username: %s, password: %s}\n", this.UserName, this.Password)
}
