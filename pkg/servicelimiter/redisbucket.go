package servicelimiter

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisBucket struct {
	client *redis.Client
}

var (
	ctx = context.Background()
)

func NewRedisBucket(addr string, password string, db int) *RedisBucket {
	if !isBucketInitialized {
		isBucketInitialized = true
		client := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		})
		return &RedisBucket{client: client}
	}
	log.Warn().Msgf("Another bucket has been initialized")
	return nil
}

func (rb *RedisBucket) refillToken(bucketConfig BucketConfig) {
	if bucketConfig.Capacity == UNLIMITED_RATE {
		return
	}

	key := bucketConfig.Key
	now := time.Now().Unix()
	lastRefillKey := key + ":last_refill"
	tokensKey := key + ":tokens"

	pipe := rb.client.TxPipeline()
	lastRefillCmd := pipe.Get(ctx, lastRefillKey)
	tokensCmd := pipe.Get(ctx, tokensKey)
	_, err := pipe.Exec(ctx)
	if err != nil && !errors.Is(redis.Nil, err) {
		return
	}

	var lastRefill int
	if errors.Is(redis.Nil, lastRefillCmd.Err()) {
		lastRefill = int(now)
		pipe.Set(ctx, lastRefillKey, now, 0)
	} else {
		lastRefill, err = strconv.Atoi(lastRefillCmd.Val())
		if err != nil {
			return
		}
	}

	var tokens int
	if errors.Is(tokensCmd.Err(), redis.Nil) {
		tokens = bucketConfig.Capacity
		pipe.Set(ctx, tokensKey, bucketConfig.Capacity, 0)
	} else {
		// convert to int
		tokens, err = strconv.Atoi(tokensCmd.Val())
		if err != nil {
			return
		}
	}

	refillRate := 1.0 / float64(bucketConfig.Ttl)
	elapsed := float64(int(now) - lastRefill)
	tokensToAdd := int(elapsed * refillRate * float64(bucketConfig.Capacity))
	updatedTokens := min(bucketConfig.Capacity, tokens+tokensToAdd)

	pipe.Set(ctx, tokensKey, updatedTokens, 0)
	pipe.Set(ctx, lastRefillKey, now, 0)
	pipe.Exec(ctx)
}

func (rb *RedisBucket) isAllow(bucketConfig BucketConfig) bool {
	tokens, err := rb.client.Get(ctx, bucketConfig.Key+":tokens").Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		return false
	}
	return tokens > 0
}

func (rb *RedisBucket) consumeToken(bucketConfig BucketConfig) bool {
	key := bucketConfig.Key
	tokensKey := key + ":tokens"
	lastRefillKey := key + ":last_refill"

	pipe := rb.client.TxPipeline()
	tokensCmd := pipe.Get(ctx, tokensKey)
	_, err := pipe.Exec(ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		return false
	}

	tokens, err := strconv.Atoi(tokensCmd.Val())
	if tokens >= 1 {
		pipe.Decr(ctx, tokensKey)
		pipe.Set(ctx, lastRefillKey, time.Now().Unix(), 0)
		pipe.Exec(ctx)
		return true
	}
	return false
}

func (rb *RedisBucket) getRemainingTokens(bucketConfig BucketConfig) int {
	tokens, err := rb.client.Get(ctx, bucketConfig.Key+":tokens").Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0
	}
	return tokens
}
