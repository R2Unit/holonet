package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init(dataSourceName string) {
	var err error
	DB, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
}
