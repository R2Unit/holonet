package netbox

import (
	"sync"
	"time"

	"github.com/holonet/core/logger"
)

type HeartbeatStatus struct {
	IsAvailable bool
	LastChecked time.Time
	LastError   string
	mutex       sync.RWMutex
}

var (
	heartbeatStatus = &HeartbeatStatus{
		IsAvailable: false,
		LastChecked: time.Time{},
		LastError:   "",
	}
)

func StartHeartbeat(client *Client) {
	logger.Info("Starting NetBox heartbeat...")

	updateHeartbeatStatus(client)

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			updateHeartbeatStatus(client)
		}
	}()
}

func updateHeartbeatStatus(client *Client) {
	isAvailable := client.IsAvailable()

	heartbeatStatus.mutex.Lock()
	defer heartbeatStatus.mutex.Unlock()

	heartbeatStatus.IsAvailable = isAvailable
	heartbeatStatus.LastChecked = time.Now()

	if !isAvailable {
		heartbeatStatus.LastError = "NetBox instance is not available"
		logger.Warn("NetBox heartbeat check failed: instance is not available")
	} else {
		heartbeatStatus.LastError = ""
		logger.Debug("NetBox heartbeat check successful")
	}
}

func GetHeartbeatStatus() HeartbeatStatus {
	heartbeatStatus.mutex.RLock()
	defer heartbeatStatus.mutex.RUnlock()

	return HeartbeatStatus{
		IsAvailable: heartbeatStatus.IsAvailable,
		LastChecked: heartbeatStatus.LastChecked,
		LastError:   heartbeatStatus.LastError,
	}
}
