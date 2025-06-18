package api

import (
	"context"
	"net/http"

	"github.com/holonet/core/logger"
)

func tokenAuthMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := GetBearerToken(r)
		if token == "" {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenInfo, err := AuthenticateToken(token, dbHandler.DB)
		if err != nil {
			logger.Error("Token authentication failed: %v", err)
			if err.Error() == "token not found or expired" || err.Error() == "token expired" {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			} else if err.Error()[:16] == "rate limit exceeded" {
				http.Error(w, err.Error(), http.StatusTooManyRequests)
				return
			} else {
				http.Error(w, "Error validating token", http.StatusInternalServerError)
				return
			}
		}

		ctx := r.Context()
		ctx = contextWithTokenInfo(ctx, tokenInfo)
		r = r.WithContext(ctx)

		handler(w, r)
	}
}

type contextKey int

const (
	tokenInfoKey contextKey = iota
)

func contextWithTokenInfo(ctx context.Context, info *TokenInfo) context.Context {
	return context.WithValue(ctx, tokenInfoKey, info)
}

func TokenInfoFromContext(ctx context.Context) (*TokenInfo, bool) {
	info, ok := ctx.Value(tokenInfoKey).(*TokenInfo)
	return info, ok
}
