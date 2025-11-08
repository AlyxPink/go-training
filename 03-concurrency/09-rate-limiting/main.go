package main

import (
	"fmt"
	"sync"
	"time"
)

// TokenBucket implements token bucket rate limiting
type TokenBucket struct {
	capacity  int
	tokens    int
	refillRate time.Duration
	mu        sync.Mutex
	lastRefill time.Time
}

// NewTokenBucket creates a token bucket rate limiter
// TODO: Implement token bucket initialization
func NewTokenBucket(capacity int, refillRate time.Duration) *TokenBucket {
	// TODO: Initialize fields
	// TODO: Start refill goroutine
	return nil
}

// Allow checks if request is allowed
// TODO: Implement token bucket algorithm
func (tb *TokenBucket) Allow() bool {
	// TODO: Refill tokens based on time passed
	// TODO: Check if token available
	// TODO: Consume token if available
	return false
}

// FixedWindowLimiter limits requests per time window
type FixedWindowLimiter struct {
	limit     int
	window    time.Duration
	count     int
	windowStart time.Time
	mu        sync.Mutex
}

// TODO: Implement FixedWindowLimiter methods

func main() {
	fmt.Println("Rate Limiting Examples")

	// Token bucket example
	limiter := NewTokenBucket(5, 1*time.Second)

	// Simulate requests
	for i := 0; i < 10; i++ {
		if limiter.Allow() {
			fmt.Println("Request", i, "allowed")
		} else {
			fmt.Println("Request", i, "rate limited")
		}
		time.Sleep(200 * time.Millisecond)
	}
}
