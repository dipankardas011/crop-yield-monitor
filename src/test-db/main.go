package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   "auth",
		Passwd: "12345",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "auth",
	}
	// Get a database handle.
	var err error
	var rows *sql.Rows
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	rows, err = db.Query("SHOW TABLES;")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	var tableName string
	for rows.Next() {
		if err := rows.Scan(&tableName); err != nil {
			panic(err)
		}
		fmt.Println(tableName)
	}

}
