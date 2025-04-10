package database

import (
	"context"
	"log"
	"time"
)

func (handler *DBHandler) StartHeartbeat() {
	heartbeatInterval := 10 * time.Second
	pingTimeout := 5 * time.Second

	for {
		ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
		err := handler.DB.PingContext(ctx)
		cancel()

		if err != nil {
			log.Printf("Heartbeat error: unable to ping the database: %v", err)
		} else {
			log.Printf("Heartbeat: Postgres at %s is responsive")
		}

		time.Sleep(heartbeatInterval)
	}
}
