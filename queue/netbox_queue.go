package queue

import (
	"sync"
	"time"

	"github.com/holonet/core/logger"
)

// NetboxQueuedRequest represents a request that has been queued due to rate limiting
type NetboxQueuedRequest struct {
	Method      string
	Endpoint    string
	Body        interface{}
	RetryAt     time.Time
	Attempts    int
	MaxAttempts int
	Result      chan NetboxQueueResult
}

// NetboxQueueResult represents the result of a queued request
type NetboxQueueResult struct {
	Data  []byte
	Error error
}

// NetboxRequestExecutor defines the interface for executing NetBox API requests
type NetboxRequestExecutor interface {
	ExecuteRequest(method, endpoint string, body interface{}) ([]byte, error)
}

// NetboxQueue manages a queue of rate-limited NetBox API requests
type NetboxQueue struct {
	executor     NetboxRequestExecutor
	queue        []*NetboxQueuedRequest
	queueMutex   sync.Mutex
	queueRunning bool
}

// NewNetboxQueue creates a new NetboxQueue instance
func NewNetboxQueue(executor NetboxRequestExecutor) *NetboxQueue {
	nq := &NetboxQueue{
		executor:     executor,
		queue:        make([]*NetboxQueuedRequest, 0),
		queueRunning: false,
	}

	// Start the queue processor
	go nq.processQueue()

	return nq
}

// QueueRequest adds a request to the queue and returns a channel for the result
func (nq *NetboxQueue) QueueRequest(method, endpoint string, body interface{}) ([]byte, error) {
	logger.Debug("Queueing NetBox request due to rate limiting: %s %s", method, endpoint)

	// Create a channel for the result
	resultChan := make(chan NetboxQueueResult, 1)

	// Create a queued request
	request := &NetboxQueuedRequest{
		Method:      method,
		Endpoint:    endpoint,
		Body:        body,
		RetryAt:     time.Now().Add(5 * time.Second), // Initial retry after 5 seconds
		Attempts:    0,
		MaxAttempts: 5, // Default max attempts
		Result:      resultChan,
	}

	// Add the request to the queue
	nq.queueMutex.Lock()
	nq.queue = append(nq.queue, request)
	nq.queueMutex.Unlock()

	// Wait for the result
	result := <-resultChan
	close(resultChan)

	return result.Data, result.Error
}

// processQueue processes queued requests when rate limits allow
func (nq *NetboxQueue) processQueue() {
	nq.queueMutex.Lock()
	nq.queueRunning = true
	nq.queueMutex.Unlock()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C

		nq.queueMutex.Lock()
		if len(nq.queue) == 0 {
			nq.queueMutex.Unlock()
			continue
		}

		now := time.Now()
		var nextRequest *NetboxQueuedRequest
		var index int

		// Find the next request that's ready to be retried
		for i, req := range nq.queue {
			if now.After(req.RetryAt) {
				nextRequest = req
				index = i
				break
			}
		}

		// If no request is ready, wait for the next tick
		if nextRequest == nil {
			nq.queueMutex.Unlock()
			continue
		}

		// Remove the request from the queue
		nq.queue = append(nq.queue[:index], nq.queue[index+1:]...)
		nq.queueMutex.Unlock()

		// Process the request
		data, err := nq.executor.ExecuteRequest(nextRequest.Method, nextRequest.Endpoint, nextRequest.Body)

		// If the request failed due to rate limiting and we haven't exceeded max attempts,
		// requeue it with a delay
		if err != nil && nextRequest.Attempts < nextRequest.MaxAttempts {
			nextRequest.Attempts++
			nextRequest.RetryAt = time.Now().Add(time.Duration(nextRequest.Attempts) * 5 * time.Second)

			nq.queueMutex.Lock()
			nq.queue = append(nq.queue, nextRequest)
			nq.queueMutex.Unlock()

			logger.Debug("Requeued NetBox request for retry (attempt %d/%d): %s %s",
				nextRequest.Attempts, nextRequest.MaxAttempts, nextRequest.Method, nextRequest.Endpoint)
		} else {
			// Send the result back through the channel
			nextRequest.Result <- NetboxQueueResult{
				Data:  data,
				Error: err,
			}

			if err != nil {
				logger.Error("Failed NetBox queued request after %d attempts: %v", nextRequest.Attempts, err)
			} else {
				logger.Debug("Successfully processed queued NetBox request: %s %s", nextRequest.Method, nextRequest.Endpoint)
			}
		}
	}
}
