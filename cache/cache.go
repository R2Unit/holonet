package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewCacheClient() (*redis.Client, error) {
	host := os.Getenv("VALKEY_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("VALKEY_PORT")
	if port == "" {
		port = "6379"
	}

	password := os.Getenv("VALKEY_PASSWORD")
	addr := fmt.Sprintf("%s:%s", host, port)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	const maxRetries = 5
	for i := 1; i <= maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := client.Ping(ctx).Err()
		cancel()

		if err == nil {
			log.Printf("Successfully connected to Redis at %s on attempt %d", addr, i)
			return client, nil
		}

		log.Printf("Attempt %d: could not connect to Redis at %s, error: %v", i, addr, err)
		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to Redis at %s after %d attempts", addr, maxRetries)
}
