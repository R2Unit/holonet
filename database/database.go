package database

import (
	"log"

	"github.com/r2unit/colours"
)

func InitializeDatabase() {
	db, err := InitDB()
	if err != nil {
		log.Fatalf(colours.Danger("Database initialization failed: %v"), err)
		return
	}
	defer db.Close()

	log.Println(colours.Success("Database connected successfully!"))

	for tableName, createQuery := range GetTables() {
		exists, err := TableExists(db, tableName)
		if err != nil {
			log.Fatalf(colours.Danger("Error checking if table %s exists: %v"), tableName, err)
		}
		if !exists {
			log.Printf(colours.Success("Creating table: %s"), tableName)
			_, err := db.Exec(createQuery)
			if err != nil {
				log.Fatalf(colours.Danger("Failed to create table %s: %v"), tableName, err)
			}
		} else {
			log.Printf(colours.Success("Table %s already exists, skipping creation."), tableName)
		}
	}
}
