package ratelimiter

import (
	"sync"
	"time"
)

type DefaultBucket struct {
	tokenBuckets    map[string]int
	lastRefillTimes map[string]time.Time
	mutex           sync.Mutex
}

func NewDefaultBucket() *DefaultBucket {
	return &DefaultBucket{
		tokenBuckets:    make(map[string]int),
		lastRefillTimes: make(map[string]time.Time),
	}
}

func (bk *DefaultBucket) refillToken(bucketConfig BucketConfig) {
	bk.mutex.Lock()
	defer bk.mutex.Unlock()

	if bucketConfig.Capacity == UNLIMITED_RATE {
		return
	}
	now := time.Now()
	// If key not in lastRefillTimes, then set lastRefillTimes to now
	if _, ok := bk.lastRefillTimes[bucketConfig.Key]; !ok {
		bk.lastRefillTimes[bucketConfig.Key] = now
	}

	// If key not in tokenBuckets, then set tokenBuckets to capacity
	if _, ok := bk.tokenBuckets[bucketConfig.Key]; !ok {
		bk.tokenBuckets[bucketConfig.Key] = bucketConfig.Capacity
	}

	refillRate := 1.0 / float64(bucketConfig.Ttl)
	elapsed := now.Sub(bk.lastRefillTimes[bucketConfig.Key]).Seconds()

	tokensToAdd := int(elapsed * refillRate * float64(bucketConfig.Capacity))
	currentToken := bk.tokenBuckets[bucketConfig.Key]
	updatedToken := currentToken + tokensToAdd
	bk.tokenBuckets[bucketConfig.Key] = min(bucketConfig.Capacity, updatedToken)
}

func (bk *DefaultBucket) isAllow(bucketConfig BucketConfig) bool {
	bk.mutex.Lock()
	defer bk.mutex.Unlock()

	return bk.tokenBuckets[bucketConfig.Key] > 0
}

func (bk *DefaultBucket) consumeToken(bucketConfig BucketConfig) bool {
	bk.mutex.Lock()
	defer bk.mutex.Unlock()

	if bk.tokenBuckets[bucketConfig.Key] >= 1 {
		bk.tokenBuckets[bucketConfig.Key] = bk.tokenBuckets[bucketConfig.Key] - 1
		bk.lastRefillTimes[bucketConfig.Key] = time.Now()
		return true
	}
	return false
}

func (bk *DefaultBucket) getRemainingTokens(bucketConfig BucketConfig) int {
	bk.mutex.Lock()
	defer bk.mutex.Unlock()

	return bk.tokenBuckets[bucketConfig.Key]
}
