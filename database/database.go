package database

import (
	"database/sql"

	"github.com/holonet/core/logger"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init(dataSourceName string) {
	var err error
	logger.Debug("Initializing database connection")
	DB, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		logger.Fatal("Error connecting to database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		logger.Fatal("Error pinging database: %v", err)
	}
	logger.Info("Database connection established successfully")
}
