package main

import (
	"log"

	"github.com/r2unit/colours"
	"github.com/r2unit/talos-core/auth"
	"github.com/r2unit/talos-core/config"
	"github.com/r2unit/talos-core/database"
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
