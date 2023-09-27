package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

const (
	DB_NAME      = "auth"
	DB_CONN_ADDR = "127.0.0.1:3306"
	dbPassword   = "12345"
	DB_USER      = "auth"
	TABLE_USERS  = "users"
	TABLE_PASS   = "passwords"
)

type DBClient struct {
	db *sql.DB
}

type AccountDBStore struct {
	userName string
	name     string
	password string
	email    string
	salt     string
}

func (this *DBClient) MySqlNewClient() error {

	// Capture connection properties.
	cfg := mysql.Config{
		User:   DB_USER,
		Passwd: dbPassword,
		Net:    "tcp",
		Addr:   DB_CONN_ADDR,
		DBName: DB_NAME,
	}
	// Get a database handle.
	var err error
	this.db, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		return err
	}
	return nil
}

func (this *DBClient) CreateUser(ac AccountSignUp) error {

	if err := allowableSizeOfUserInputs(ac); err != nil {
		return err
	}

	// there should be no user with same username
	if account, err := this.GetUserByUsername(ac.UserName); err != nil {
		if err == sql.ErrNoRows {
			log.Println("No user found")
		} else {
			return err
		}
	} else if account != nil {
		// another user is there
		return fmt.Errorf("UserName already exist")
	}

	result, err := this.db.Exec("INSERT INTO users (username, name, email) VALUES (?, ?, ?)", ac.UserName, ac.Name, ac.Email)
	if err != nil {
		return fmt.Errorf("add details to %s table: %v", TABLE_USERS, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("add details to %s table: %v", TABLE_USERS, err)
	}

	log.Printf("Add details to the `%s` table for username: %s the sql result: %v", TABLE_USERS, ac.UserName, id)

	////////  PASSWORDS

	hashed, salt, err := GenerateHashForPasswordAndSalt(ac.Password)
	if err != nil {
		return err
	}

	result, err = this.db.Exec("INSERT INTO passwords (username, password, salt) VALUES (?, ?, ?)", ac.UserName, hashed, salt)
	if err != nil {
		return fmt.Errorf("add details to %s table: %v", TABLE_PASS, err)
	}

	id, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("add details to %s table: %v", TABLE_PASS, err)
	}

	log.Printf("Add details to the `%s` table for username: %s the sql result: %v", TABLE_PASS, ac.UserName, id)

	return nil
}

func (this *DBClient) DeleteUser(username string) error {
	// cleanup the users table and passwords table
	return nil
}

func (this *DBClient) ListOfUsers() ([]AccountDBStore, error) {
	return nil, nil
}

func (this *DBClient) GetPasswordByUsername(username string) (*AccountDBStore, error) {
	// An album to hold data from the returned row.
	user := AccountDBStore{}

	row := this.db.QueryRow("SELECT * FROM passwords WHERE username = ?", username)

	if err := row.Scan(&user.userName, &user.password, &user.salt); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("getPasswordByUsername %s: %v", username, err)
	}
	return &user, nil
}

func (this *DBClient) GetUserByUsername(username string) (*AccountDBStore, error) {
	// An album to hold data from the returned row.
	user := AccountDBStore{}

	row := this.db.QueryRow("SELECT * FROM users WHERE username = ?", username)

	if err := row.Scan(&user.userName, &user.name, &user.email); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("getUserByUsername %s: %v", username, err)
	}
	return &user, nil
}
