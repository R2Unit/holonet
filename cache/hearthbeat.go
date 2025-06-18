package cache

import (
	"context"
	"github.com/holonet/core/logger"
	"time"
)

func StartHeartbeat(client CacheClient) {
	addr := client.Options().Addr
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := client.Ping(ctx)
		cancel()

		if err != nil {
			logger.Error("Heartbeat error: unable to ping Redis at %s: %v", addr, err)
		} else {
			logger.Debug("Heartbeat: Redis at %s is responsive", addr)
		}

		time.Sleep(10 * time.Second)
	}
}
