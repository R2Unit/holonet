package main

import (
	"github.com/holonet/core/database"
	_ "github.com/holonet/core/database/tables"
	"log"
)

func main() {
	dbHandler, err := database.NewDBHandler()
	if err != nil {
		log.Fatalf("Failed to initialize DB: %v", err)
	}

	log.Printf("Registered %d table(s) for migration.", database.RegisteredTableCount())

	if err := dbHandler.Migrate(); err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	log.Println("Database migrations completed successfully.")
}
