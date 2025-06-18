package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/holonet/core/logger"
)

type TokenInfo struct {
	ID                 int
	UserID             int
	Token              string
	PolicyID           sql.NullInt32
	RateLimitPerMin    int
	MaxRequestsPerDay  int
	RequestCount       int
	LastRequestAt      sql.NullTime
	RequestsToday      int
	RequestsTodayReset sql.NullTime
	ExpiresAt          time.Time
}

type TokenCache struct {
	tokens map[string]*TokenInfo
	mutex  sync.RWMutex
}

var tokenCache = &TokenCache{
	tokens: make(map[string]*TokenInfo),
}

func (c *TokenCache) Get(token string) (*TokenInfo, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	info, exists := c.tokens[token]
	return info, exists
}

func (c *TokenCache) Set(token string, info *TokenInfo) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.tokens[token] = info
}

func (c *TokenCache) Remove(token string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.tokens, token)
}

func GetTokenInfo(token string, db *sql.DB) (*TokenInfo, error) {
	if info, exists := tokenCache.Get(token); exists {
		if info.ExpiresAt.Before(time.Now()) {
			tokenCache.Remove(token)
			return nil, fmt.Errorf("token expired")
		}
		return info, nil
	}

	query := `
		SELECT t.id, t.user_id, t.token, t.policy_id, 
		       COALESCE(p.rate_limit_per_min, 60) as rate_limit_per_min, 
		       COALESCE(p.max_requests_per_day, 1000) as max_requests_per_day,
		       t.request_count, t.last_request_at, 
		       t.requests_today, t.requests_today_reset, t.expires_at
		FROM tokens t
		LEFT JOIN token_policies p ON t.policy_id = p.id
		WHERE t.token = $1 AND t.expires_at > NOW()
	`

	info := &TokenInfo{}
	err := db.QueryRow(query, token).Scan(
		&info.ID,
		&info.UserID,
		&info.Token,
		&info.PolicyID,
		&info.RateLimitPerMin,
		&info.MaxRequestsPerDay,
		&info.RequestCount,
		&info.LastRequestAt,
		&info.RequestsToday,
		&info.RequestsTodayReset,
		&info.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("token not found or expired")
		}
		return nil, fmt.Errorf("error querying token: %w", err)
	}

	tokenCache.Set(token, info)
	return info, nil
}

func UpdateTokenUsage(info *TokenInfo, db *sql.DB) error {
	now := time.Now()
	info.RequestCount++
	info.LastRequestAt = sql.NullTime{Time: now, Valid: true}

	if !info.RequestsTodayReset.Valid || now.After(info.RequestsTodayReset.Time) {
		info.RequestsToday = 1
		midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		info.RequestsTodayReset = sql.NullTime{Time: midnight, Valid: true}
	} else {
		info.RequestsToday++
	}

	query := `
		UPDATE tokens
		SET request_count = $1, 
		    last_request_at = $2, 
		    requests_today = $3, 
		    requests_today_reset = $4,
		    updated_at = NOW()
		WHERE id = $5
	`

	_, err := db.Exec(
		query,
		info.RequestCount,
		info.LastRequestAt,
		info.RequestsToday,
		info.RequestsTodayReset,
		info.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating token usage: %w", err)
	}

	tokenCache.Set(info.Token, info)
	return nil
}

func CheckRateLimit(info *TokenInfo) (bool, string) {
	now := time.Now()

	if info.RequestsToday >= info.MaxRequestsPerDay {
		nextReset := info.RequestsTodayReset.Time
		return false, fmt.Sprintf("Daily request limit exceeded. Limit resets at %s", nextReset.Format(time.RFC3339))
	}

	if info.LastRequestAt.Valid {
		lastRequest := info.LastRequestAt.Time
		if now.Sub(lastRequest) < time.Minute {
			if info.RequestsToday > info.RateLimitPerMin {
				return false, fmt.Sprintf("Rate limit exceeded. Maximum %d requests per minute allowed.", info.RateLimitPerMin)
			}
		}
	}

	return true, ""
}

func GetBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}

func AuthenticateToken(token string, db *sql.DB) (*TokenInfo, error) {
	info, err := GetTokenInfo(token, db)
	if err != nil {
		return nil, err
	}

	allowed, message := CheckRateLimit(info)
	if !allowed {
		return nil, fmt.Errorf("rate limit exceeded: %s", message)
	}

	err = UpdateTokenUsage(info, db)
	if err != nil {
		logger.Error("Failed to update token usage: %v", err)
	}

	return info, nil
}
