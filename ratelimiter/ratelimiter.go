package ratelimiter

import (
	"sync"
	"time"
)

type Bucket struct {
	tokens     int
	capacity   int
	refillRate int
	mu         sync.Mutex
}

func NewBucket(capacity, refillRate int) *Bucket {
	return &Bucket{tokens: capacity, capacity: capacity, refillRate: refillRate}
}

func (b *Bucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.tokens > 0 {
		b.tokens--
		return true
	}
	return false
}

func (b *Bucket) Refill() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.tokens < b.capacity {
		b.tokens++
	}
}

type RateLimiter struct {
	buckets map[string]*Bucket
	mu      sync.RWMutex
}

func NewRateLimiter(defaultCapacity, defaultRefill int) *RateLimiter {
	rl := &RateLimiter{buckets: make(map[string]*Bucket)}
	go rl.refillAll(defaultRefill)
	return rl
}

func (r *RateLimiter) refillAll(refillRate int) {
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		r.mu.RLock()
		for _, bucket := range r.buckets {
			bucket.Refill()
		}
		r.mu.RUnlock()
	}
}

func (r *RateLimiter) Allow(ip string) bool {
	r.mu.RLock()
	bucket, ok := r.buckets[ip]
	r.mu.RUnlock()

	if !ok {
		return false
	}
	return bucket.Allow()
}

func (r *RateLimiter) SetClientLimit(ip string, capacity, refill int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.buckets[ip] = NewBucket(capacity, refill)
}
