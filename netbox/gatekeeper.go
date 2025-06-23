package netbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/holonet/core/logger"
	"github.com/holonet/core/queue"
)

type Gatekeeper struct {
	client       *Client
	rateLimiter  *RateLimiter
	cacheEnabled bool
	cache        map[string]CachedResponse
	cacheMutex   sync.RWMutex
	cacheExpiry  time.Duration

	netboxQueue *queue.NetboxQueue
}

type CachedResponse struct {
	Data      []byte
	Timestamp time.Time
}

func (g *Gatekeeper) ExecuteRequest(method, endpoint string, body interface{}) ([]byte, error) {
	return g.executeRequestDirect(method, endpoint, body)
}

type RateLimiter struct {
	requestsPerMinute int
	requestCount      int
	resetTime         time.Time
	mutex             sync.Mutex
}

func NewGatekeeper(client *Client) *Gatekeeper {
	gk := &Gatekeeper{
		client: client,
		rateLimiter: &RateLimiter{
			requestsPerMinute: 100,
			resetTime:         time.Now().Add(time.Minute),
		},
		cacheEnabled: true,
		cache:        make(map[string]CachedResponse),
		cacheExpiry:  5 * time.Minute,
	}

	gk.netboxQueue = queue.NewNetboxQueue(gk)

	return gk
}

func (g *Gatekeeper) SetRateLimit(requestsPerMinute int) {
	g.rateLimiter.mutex.Lock()
	defer g.rateLimiter.mutex.Unlock()
	g.rateLimiter.requestsPerMinute = requestsPerMinute
}

func (g *Gatekeeper) SetCacheEnabled(enabled bool) {
	g.cacheEnabled = enabled
}

func (g *Gatekeeper) SetCacheExpiry(duration time.Duration) {
	g.cacheExpiry = duration
}

func (g *Gatekeeper) Request(method, endpoint string, body interface{}) ([]byte, error) {
	cacheKey := fmt.Sprintf("%s:%s", method, endpoint)
	if method == http.MethodGet && g.cacheEnabled {
		g.cacheMutex.RLock()
		if cachedResp, ok := g.cache[cacheKey]; ok {
			if time.Since(cachedResp.Timestamp) < g.cacheExpiry {
				g.cacheMutex.RUnlock()
				return cachedResp.Data, nil
			}
		}
		g.cacheMutex.RUnlock()
	}

	if !g.checkRateLimit() {
		return g.netboxQueue.QueueRequest(method, endpoint, body)
	}

	return g.executeRequestDirect(method, endpoint, body)
}

func (g *Gatekeeper) executeRequestDirect(method, endpoint string, body interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s/api/%s", g.client.Host, endpoint)
	cacheKey := fmt.Sprintf("%s:%s", method, endpoint)

	var req *http.Request
	var err error

	if body != nil {
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
		bodyReader := bytes.NewReader(bodyJSON)
		req, err = http.NewRequest(method, url, bodyReader)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}
	}

	if g.client.Token != "" {
		req.Header.Set("Authorization", "Token "+g.client.Token)
	}

	resp, err := g.client.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(errorBody))
	}

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if method == http.MethodGet && g.cacheEnabled {
		g.cacheMutex.Lock()
		g.cache[cacheKey] = CachedResponse{
			Data:      responseData,
			Timestamp: time.Now(),
		}
		g.cacheMutex.Unlock()
	}

	return responseData, nil
}

func (g *Gatekeeper) checkRateLimit() bool {
	g.rateLimiter.mutex.Lock()
	defer g.rateLimiter.mutex.Unlock()

	now := time.Now()
	if now.After(g.rateLimiter.resetTime) {
		g.rateLimiter.requestCount = 0
		g.rateLimiter.resetTime = now.Add(time.Minute)
	}

	if g.rateLimiter.requestCount >= g.rateLimiter.requestsPerMinute {
		return false
	}

	g.rateLimiter.requestCount++
	return true
}

func (g *Gatekeeper) ClearCache() {
	g.cacheMutex.Lock()
	g.cache = make(map[string]CachedResponse)
	g.cacheMutex.Unlock()
	logger.Info("NetBox API response cache cleared")
}
