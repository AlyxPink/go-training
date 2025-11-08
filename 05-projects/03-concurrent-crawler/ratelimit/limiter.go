package ratelimit

import (
	"context"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiter *rate.Limiter
}

func NewRateLimiter(requestsPerSecond float64) *RateLimiter {
	// TODO: Create token bucket limiter
	return &RateLimiter{
		limiter: rate.NewLimiter(rate.Limit(requestsPerSecond), int(requestsPerSecond)),
	}
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
	// TODO: Wait for token
	return rl.limiter.Wait(ctx)
}

func (rl *RateLimiter) Allow() bool {
	// TODO: Check if request allowed without blocking
	return rl.limiter.Allow()
}
