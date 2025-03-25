package context

import (
	"context"
	"time"
)

type key int

const (
	// RequestIDKey is the key for request ID in context
	RequestIDKey key = iota
	// UserIDKey is the key for user ID in context
	UserIDKey
)

// GetRequestID retrieves the request ID from context
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// WithRequestID adds request ID to context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// GetUserID retrieves the user ID from context
func GetUserID(ctx context.Context) int64 {
	if ctx == nil {
		return 0
	}
	if userID, ok := ctx.Value(UserIDKey).(int64); ok {
		return userID
	}
	return 0
}

// WithUserID adds user ID to context
func WithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func WithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, timeout)
}
