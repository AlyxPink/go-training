package main

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucket struct {
	capacity   int
	tokens     int
	refillRate time.Duration
	mu         sync.Mutex
	lastRefill time.Time
}

func NewTokenBucket(capacity int, refillRate time.Duration) *TokenBucket {
	tb := &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}

	// Background refill
	go func() {
		ticker := time.NewTicker(refillRate)
		defer ticker.Stop()

		for range ticker.C {
			tb.mu.Lock()
			if tb.tokens < tb.capacity {
				tb.tokens++
			}
			tb.mu.Unlock()
		}
	}()

	return tb
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

type FixedWindowLimiter struct {
	limit       int
	window      time.Duration
	count       int
	windowStart time.Time
	mu          sync.Mutex
}

func NewFixedWindowLimiter(limit int, window time.Duration) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		limit:       limit,
		window:      window,
		count:       0,
		windowStart: time.Now(),
	}
}

func (f *FixedWindowLimiter) Allow() bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := time.Now()

	// Reset window if expired
	if now.Sub(f.windowStart) >= f.window {
		f.count = 0
		f.windowStart = now
	}

	if f.count < f.limit {
		f.count++
		return true
	}
	return false
}

func main() {
	fmt.Println("Rate Limiting Examples")

	limiter := NewTokenBucket(5, 200*time.Millisecond)

	for i := 0; i < 10; i++ {
		if limiter.Allow() {
			fmt.Println("Request", i, "allowed")
		} else {
			fmt.Println("Request", i, "rate limited")
		}
		time.Sleep(100 * time.Millisecond)
	}
}
