package cache

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type ValkeyCacheClient struct {
	addr     string
	password string
	db       int
	conn     net.Conn
}

func NewValkeyCacheClient() (*ValkeyCacheClient, error) {
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

	client := &ValkeyCacheClient{
		addr:     addr,
		password: password,
		db:       0,
	}

	const maxRetries = 5
	for i := 1; i <= maxRetries; i++ {
		err := client.connect()
		if err == nil {
			log.Printf("Successfully connected to Redis at %s on attempt %d", addr, i)
			return client, nil
		}
		log.Printf("Attempt %d: could not connect to Redis at %s, error: %v", i, addr, err)
		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to Redis at %s after %d attempts", addr, maxRetries)
}

func (c *ValkeyCacheClient) connect() error {
	var err error
	c.conn, err = net.DialTimeout("tcp", c.addr, 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	if c.password != "" {
		_, err = c.executeCommand(fmt.Sprintf("AUTH %s\r\n", c.password))
		if err != nil {
			c.conn.Close()
			return fmt.Errorf("failed to authenticate: %w", err)
		}
	}

	if c.db != 0 {
		_, err = c.executeCommand(fmt.Sprintf("SELECT %d\r\n", c.db))
		if err != nil {
			c.conn.Close()
			return fmt.Errorf("failed to select database: %w", err)
		}
	}

	return nil
}

func (c *ValkeyCacheClient) executeCommand(cmd string) (string, error) {
	if c.conn == nil {
		return "", fmt.Errorf("not connected to Redis")
	}
	err := c.conn.SetDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return "", fmt.Errorf("failed to set deadline: %w", err)
	}
	_, err = c.conn.Write([]byte(cmd))
	if err != nil {
		return "", fmt.Errorf("failed to send command: %w", err)
	}
	buf := make([]byte, 1024)
	n, err := c.conn.Read(buf)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}
	return string(buf[:n]), nil
}

func (c *ValkeyCacheClient) Ping(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	resultCh := make(chan error, 1)
	go func() {
		_, err := c.executeCommand("PING\r\n")
		resultCh <- err
	}()
	select {
	case err := <-resultCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *ValkeyCacheClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *ValkeyCacheClient) Options() *ValkeyOptions {
	return &ValkeyOptions{
		Addr:     c.addr,
		Password: c.password,
		DB:       c.db,
	}
}

type ValkeyOptions struct {
	Addr     string
	Password string
	DB       int
}
