package cache

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestValkeyCacheClient(t *testing.T) {
	if os.Getenv("SKIP_REDIS_TESTS") == "true" {
		t.Skip("Skipping Redis tests")
	}
	client, err := NewValkeyCacheClient()
	if err != nil {
		t.Fatalf("Failed to create valkey cache client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx)
	if err != nil {
		t.Fatalf("Failed to ping Redis: %v", err)
	}

	options := client.Options()
	if options == nil {
		t.Fatal("Options method returned nil")
	}
	if options.Addr == "" {
		t.Fatal("Options.Addr is empty")
	}
}

func TestNewCacheClient(t *testing.T) {
	if os.Getenv("SKIP_REDIS_TESTS") == "true" {
		t.Skip("Skipping Redis tests")
	}
	client, err := NewCacheClient()
	if err != nil {
		t.Fatalf("Failed to create cache client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx)
	if err != nil {
		t.Fatalf("Failed to ping Redis: %v", err)
	}

	options := client.Options()
	if options == nil {
		t.Fatal("Options method returned nil")
	}
	if options.Addr == "" {
		t.Fatal("Options.Addr is empty")
	}
}
