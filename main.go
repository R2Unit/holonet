package main

import (
	"log"

	"github.com/quanza/talos-core/auth"
	"github.com/quanza/talos-core/config"
	"github.com/quanza/talos-core/database"
	"github.com/r2unit/colours"
)

func main() {
	cfg := config.LoadConfig()

	// Environment Debug
	log.Printf(colours.Debug("Settings Config: %+v"), cfg.Settings)
	log.Printf(colours.Debug("Postgres Config: %+v"), cfg.Postgres)
	log.Printf(colours.Debug("S3 Config: %+v"), cfg.S3)

	database.InitializeDatabase()

	// Database Debug
	log.Printf(colours.Debug("Connecting to DB: host=%s port=%d user=%s dbname=%s"),
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.DBName)

	// Initialize Authentication
	db, err := auth.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize authentication:", err)
	}
	defer db.Close()

	// Initialize Database
}
