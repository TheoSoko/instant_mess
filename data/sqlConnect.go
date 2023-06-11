package data

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

// Global db struct
var db *sql.DB

var mysqlConf = mysql.Config{
	User:   "root",
	Passwd: "",
	Net:    "tcp",
	Addr:   "127.0.0.1:3306",
	DBName: "zemus_api",
	AllowNativePasswords: true,
}

func SqlConnect() {
	var err error

	db, err = sql.Open("mysql", mysqlConf.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected to db")
}
