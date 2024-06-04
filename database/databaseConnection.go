package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", "root", "root", "127.0.0.1", "3306", "marketplace")
	log.Printf("Connecting to database with DSN: %s", dsn)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

func DBinstance() *sql.DB {
	return db
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}
