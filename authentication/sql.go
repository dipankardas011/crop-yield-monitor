package main

import (
	"database/sql"
	"errors"
	"fmt"
)

func testMysqlLoginAndSignup() {
	mysql := &DBClient{}

	if err := mysql.MySqlNewClient(); err != nil {
		panic(err)
	}
	signup := AccountSignUp{
		Name:     "Dipankar Dsa",
		UserName: "dipankar",
		Password: "1234",
		Email:    "20051554@kiit.ac.in",
	}

	if err := mysql.CreateUser(signup); err != nil {
		panic(err)
	}

	abcd, err := mysql.GetPasswordByUsername(signup.UserName)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	err = func() error {
		getHash := genHash(signup.Password + abcd.salt)
		if abcd.password != getHash {
			return errors.New("invalid password")
		}
		return nil
	}()

	if err != nil {
		panic(err)
	}
	fmt.Println("signin success")

}
