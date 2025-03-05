package ratelimit

import (
	"sync"
	"time"
)

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	mu           sync.Mutex
	tokens       int
	maxTokens    int
	refillRate   int
	lastRefill   time.Time
	refillPeriod time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxTokens int, refillRate int, refillPeriod time.Duration) *RateLimiter {
	return &RateLimiter{
		tokens:       maxTokens,
		maxTokens:    maxTokens,
		refillRate:   refillRate,
		lastRefill:   time.Now(),
		refillPeriod: refillPeriod,
	}
}

// Acquire attempts to acquire a token
func (r *RateLimiter) Acquire() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Refill tokens if needed
	now := time.Now()
	elapsed := now.Sub(r.lastRefill)
	if elapsed >= r.refillPeriod {
		refillCount := int(elapsed / r.refillPeriod)
		r.tokens = min(r.maxTokens, r.tokens+refillCount*r.refillRate)
		r.lastRefill = now
	}

	if r.tokens > 0 {
		r.tokens--
		return true
	}

	return false
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
