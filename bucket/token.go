package bucket

import (
	"sync"
	"time"
)

// TokenBucket is a token bucket.
type TokenBucket struct {
	rate           int
	maxTokens      int
	currentTokens  int
	lastRefillTime time.Time
	lock           sync.Mutex
	timer          timer
}

// NewTokenBucket returns new TokenBucket.
func NewTokenBucket(rate int, maxTokens int) *TokenBucket {
	return &TokenBucket{
		rate:           rate,
		maxTokens:      maxTokens,
		currentTokens:  maxTokens,
		lastRefillTime: time.Now(),
		lock:           sync.Mutex{},
		timer:          standardTimer{},
	}
}

// Fill refills bucket and adds tokens to the bucket.
// Returns the amount of tokens beyond the edges of the bucket.
func (b *TokenBucket) Fill(tokens int) (beyond int) {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.refill()
	b.currentTokens -= tokens
	if b.currentTokens < 0 {
		return -b.currentTokens
	}
	return 0
}

func (b *TokenBucket) refill() {
	now := b.timer.Now()
	end := b.timer.Since(b.lastRefillTime)
	toAddTokens := (end.Nanoseconds() * int64(b.rate)) / 1_000_000_000
	b.currentTokens = int(min(float64(int64(b.currentTokens)+toAddTokens), float64(b.maxTokens)))
	b.lastRefillTime = now
}

// We need it, so we can test TokenBucket.
type timer interface {
	Now() time.Time
	Since(time.Time) time.Duration
}

type standardTimer struct{}

func (standardTimer) Now() time.Time                    { return time.Now() }
func (s standardTimer) Since(t time.Time) time.Duration { return time.Since(t) }
