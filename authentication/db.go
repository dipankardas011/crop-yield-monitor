package main

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

const (
	DB_NAME      = "auth"
	DB_CONN_ADDR = "127.0.0.1:3306"
	dbPassword   = "12345"
	DB_USER      = "auth"
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

func (this *DBClient) mysqlNewClient() error {

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

func (this *DBClient) createUser(ac AccountSignUp) error {
	return nil
}

func (this *DBClient) deleteUser(username string) error {
	return nil
}

func (this *DBClient) listOfUsers() ([]AccountDBStore, error) {
	return nil, nil
}
