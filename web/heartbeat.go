package web

import (
	"log"
	"net"
	"time"
)

const heartbeatInterval = 30 * time.Second

func StartHeartbeat(conn net.Conn, stopCh chan struct{}) {
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			if err := writeFrame(conn, 9, []byte("ping")); err != nil {
				log.Printf("Error sending ping: %v", err)
				return
			}
			log.Println("Sent heartbeat ping")
		}
	}
}
