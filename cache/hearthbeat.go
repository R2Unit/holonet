package cache

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func StartHeartbeat(client *redis.Client) {
	addr := client.Options().Addr
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := client.Ping(ctx).Err()
		cancel()

		if err != nil {
			log.Printf("Heartbeat error: unable to ping Redis at %s: %v", addr, err)
		} else {
			log.Printf("Heartbeat: Redis at %s is responsive", addr)
		}

		time.Sleep(10 * time.Second)
	}
}
