package data

import (
	"database/sql"
	"fmt"
	"log"

	"os"

	"github.com/go-sql-driver/mysql"
)

// Global db struct
var db *sql.DB

func SqlConn() error {
	var err error

	var mysqlConf = mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASSWORD"),
		Net:                  os.Getenv("DB_NETWORK"),
		Addr:                 os.Getenv("DB_ADDRESS"),
		DBName:               os.Getenv("DB_NAME"),
		AllowNativePasswords: true,
	}

	db, err = sql.Open("mysql", mysqlConf.FormatDSN())
	if err != nil {
		log.Println("Erreur de co à la bdd")
		return err
	}

	err = db.Ping()
	if err != nil {
		log.Println("Erreur de ping à la bdd")
		return err
	}

	fmt.Println("Connected to db")

	return nil
}

func WhatIsdb() {
	fmt.Println("whatIsdb : ", db)
}
