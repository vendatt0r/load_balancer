package limiter

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type TokenBucket struct {
	capacity     int
	tokens       int
	refillRate   int
	lastRefilled time.Time
	mu           sync.Mutex
}

type Limiter struct {
	buckets    map[string]*TokenBucket
	mu         sync.Mutex
	capacity   int
	refillRate int
}

func NewLimiter(capacity, refillRate int) *Limiter {
	lim := &Limiter{
		buckets:    make(map[string]*TokenBucket),
		capacity:   capacity,
		refillRate: refillRate,
	}
	go lim.refillBuckets()
	return lim
}

func (l *Limiter) getBucket(clientID string) *TokenBucket {
	l.mu.Lock()
	defer l.mu.Unlock()
	bucket, exists := l.buckets[clientID]
	if !exists {
		bucket = &TokenBucket{
			capacity:     l.capacity,
			tokens:       l.capacity,
			refillRate:   l.refillRate,
			lastRefilled: time.Now(),
		}
		l.buckets[clientID] = bucket
	}
	return bucket
}

func (l *Limiter) Allow(clientID string) bool {
	bucket := l.getBucket(clientID)
	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(bucket.lastRefilled).Seconds()
	newTokens := int(elapsed * float64(bucket.refillRate))

	if newTokens > 0 {
		bucket.tokens += newTokens
		if bucket.tokens > bucket.capacity {
			bucket.tokens = bucket.capacity
		}
		bucket.lastRefilled = now
	}

	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}
	return false
}

func (l *Limiter) refillBuckets() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		l.mu.Lock()
		for _, bucket := range l.buckets {
			bucket.mu.Lock()
			bucket.tokens += bucket.refillRate
			if bucket.tokens > bucket.capacity {
				bucket.tokens = bucket.capacity
			}
			bucket.mu.Unlock()
		}
		l.mu.Unlock()
	}
}

func GetClientIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
