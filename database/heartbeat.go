package database

import (
	"context"
	"github.com/holonet/core/logger"
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
			logger.Error("Heartbeat error: unable to ping the database: %v", err)
		} else {
			logger.Debug("Heartbeat: Postgres at %s is responsive")
		}

		time.Sleep(heartbeatInterval)
	}
}
