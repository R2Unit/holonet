package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// DB represents a shared database connection pool to be used across the application for executing SQL queries.
var DB *sql.DB

// Init initializes the database connection using the provided data source name and verifies the connection.
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
