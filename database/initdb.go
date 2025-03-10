package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/r2unit/colours"
	"github.com/r2unit/holonet/config"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	cfg := config.LoadConfig()

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf(colours.Danger("Failed to open database connection: %v"), err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Fatalf(colours.Danger("Database connection failed: %v"), err)
		return nil, err
	}

	log.Println(colours.Success("Database connection established successfully"))
	return db, nil
}
