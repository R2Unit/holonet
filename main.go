package main

import (
	"log"
	"net/http"
	"time"

	"github.com/r2unit/holonet/api"
	"github.com/r2unit/holonet/auth"
	"github.com/r2unit/holonet/config"
	"github.com/r2unit/holonet/controller"
	"github.com/r2unit/holonet/database"

	"github.com/r2unit/go-colours"
)

func main() {
	// Initialization van Environmonts...
	// TODO: Kijken om een wait or retry check in te bouwen in de Postgres connectivity. Canceled terug als de database neit up us.
	log.Println(colours.Success("Loading Environments..."))
	time.Sleep(3 * time.Second)
	cfg := config.LoadConfig()
	time.Sleep(2 * time.Second)
	log.Printf(colours.Success("Settings Config: %+v"), cfg.Settings)
	time.Sleep(1 * time.Second)
	log.Printf(colours.Success("Postgres Config: %+v"), cfg.Postgres)
	time.Sleep(1 * time.Second)
	log.Printf(colours.Success("S3 Config: %+v"), cfg.S3)
	time.Sleep(2 * time.Second)
	log.Println(colours.Success("Environments loaded =)"))

	// Initialization of Database
	database.InitializeDatabase()

	log.Printf(colours.Success("Connecting to DB: host=%s port=%d user=%s dbname=%s"),
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.DBName)

	// Initialize Authentication
	db, err := auth.InitDB()
	if err != nil {
		log.Fatal(colours.Error("Failed to initialize authentication:"), err)
	}
	defer db.Close()

	controller.SetDB(db)

	// http.HandleFunc("/api/ussr"),
	http.HandleFunc("/api/task", api.TaskHandler)
	http.HandleFunc("/ws", controller.WSHandler)
	http.HandleFunc("/workers", api.WorkersHandler)
	http.HandleFunc("/queue", api.QueueHandler)

	log.Println("Core server starting on :8080")
	if err := http.ListenAndServe(":443", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
