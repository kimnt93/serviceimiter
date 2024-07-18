package ratelimiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	Buckets     map[string]int
	LastRefills map[string]time.Time
}

var (
	tokenBucket *TokenBucket
	mutex       sync.Mutex
)

func NewDefaultTokenBucket() {
	tokenBucket = &TokenBucket{
		Buckets:     make(map[string]int),
		LastRefills: make(map[string]time.Time),
	}
}

func refillToken(bucketConfig BucketConfig) {
	mutex.Lock()
	defer mutex.Unlock()

	if bucketConfig.Capacity == UNLIMITED_RATE {
		return
	}
	now := time.Now()
	// If key not in LastRefills, then set LastRefills to now
	if _, ok := tokenBucket.LastRefills[bucketConfig.Key]; !ok {
		tokenBucket.LastRefills[bucketConfig.Key] = now
	}

	// If key not in Buckets, then set Buckets to capacity
	if _, ok := tokenBucket.Buckets[bucketConfig.Key]; !ok {
		tokenBucket.Buckets[bucketConfig.Key] = bucketConfig.Capacity
	}

	refillRate := 1.0 / float64(bucketConfig.Ttl)
	elapsed := now.Sub(tokenBucket.LastRefills[bucketConfig.Key]).Seconds()

	tokensToAdd := int(elapsed * refillRate * float64(bucketConfig.Capacity))
	currentToken := tokenBucket.Buckets[bucketConfig.Key]
	updatedToken := currentToken + tokensToAdd
	tokenBucket.Buckets[bucketConfig.Key] = min(bucketConfig.Capacity, updatedToken)
	tokenBucket.LastRefills[bucketConfig.Key] = now
}

func isAllow(bucketConfig BucketConfig) bool {
	mutex.Lock()
	defer mutex.Unlock()

	return tokenBucket.Buckets[bucketConfig.Key] > 0
}

func consumeToken(bucketConfig BucketConfig) bool {
	mutex.Lock()
	defer mutex.Unlock()

	if tokenBucket.Buckets[bucketConfig.Key] >= 1 {
		tokenBucket.Buckets[bucketConfig.Key] = tokenBucket.Buckets[bucketConfig.Key] - 1
		tokenBucket.LastRefills[bucketConfig.Key] = time.Now()
		return true
	}
	return false
}

func getRemainingTokens(bucketConfig BucketConfig) int {
	mutex.Lock()
	defer mutex.Unlock()

	return tokenBucket.Buckets[bucketConfig.Key]
}
