package queue

import (
	"sync"
	"time"

	"github.com/holonet/core/logger"
)

type NetboxQueuedRequest struct {
	Method      string
	Endpoint    string
	Body        interface{}
	RetryAt     time.Time
	Attempts    int
	MaxAttempts int
	Result      chan NetboxQueueResult
}

type NetboxQueueResult struct {
	Data  []byte
	Error error
}

type NetboxRequestExecutor interface {
	ExecuteRequest(method, endpoint string, body interface{}) ([]byte, error)
}

type NetboxQueue struct {
	executor     NetboxRequestExecutor
	queue        []*NetboxQueuedRequest
	queueMutex   sync.Mutex
	queueRunning bool
}

func NewNetboxQueue(executor NetboxRequestExecutor) *NetboxQueue {
	nq := &NetboxQueue{
		executor:     executor,
		queue:        make([]*NetboxQueuedRequest, 0),
		queueRunning: false,
	}

	go nq.processQueue()

	return nq
}

func (nq *NetboxQueue) QueueRequest(method, endpoint string, body interface{}) ([]byte, error) {
	logger.Debug("Queueing NetBox request due to rate limiting: %s %s", method, endpoint)

	resultChan := make(chan NetboxQueueResult, 1)

	request := &NetboxQueuedRequest{
		Method:      method,
		Endpoint:    endpoint,
		Body:        body,
		RetryAt:     time.Now().Add(5 * time.Second),
		Attempts:    0,
		MaxAttempts: 5,
		Result:      resultChan,
	}

	nq.queueMutex.Lock()
	nq.queue = append(nq.queue, request)
	nq.queueMutex.Unlock()

	result := <-resultChan
	close(resultChan)

	return result.Data, result.Error
}

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

		for i, req := range nq.queue {
			if now.After(req.RetryAt) {
				nextRequest = req
				index = i
				break
			}
		}

		if nextRequest == nil {
			nq.queueMutex.Unlock()
			continue
		}

		nq.queue = append(nq.queue[:index], nq.queue[index+1:]...)
		nq.queueMutex.Unlock()

		data, err := nq.executor.ExecuteRequest(nextRequest.Method, nextRequest.Endpoint, nextRequest.Body)

		if err != nil && nextRequest.Attempts < nextRequest.MaxAttempts {
			nextRequest.Attempts++
			nextRequest.RetryAt = time.Now().Add(time.Duration(nextRequest.Attempts) * 5 * time.Second)

			nq.queueMutex.Lock()
			nq.queue = append(nq.queue, nextRequest)
			nq.queueMutex.Unlock()

			logger.Debug("Requeued NetBox request for retry (attempt %d/%d): %s %s",
				nextRequest.Attempts, nextRequest.MaxAttempts, nextRequest.Method, nextRequest.Endpoint)
		} else {
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
