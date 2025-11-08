package ratelimit

import (
	"context"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiter *rate.Limiter
}

func NewRateLimiter(requestsPerSecond float64) *RateLimiter {
	// Create token bucket limiter with burst equal to rate
	return &RateLimiter{
		limiter: rate.NewLimiter(rate.Limit(requestsPerSecond), int(requestsPerSecond)),
	}
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
	// Wait for token, respecting context cancellation
	return rl.limiter.Wait(ctx)
}

func (rl *RateLimiter) Allow() bool {
	// Check if request allowed without blocking
	return rl.limiter.Allow()
}
