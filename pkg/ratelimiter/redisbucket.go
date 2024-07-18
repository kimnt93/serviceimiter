package ratelimiter

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type RedisBucket struct {
	client *redis.Client
}

var (
	Ctx = context.Background()
)

func NewRedisBucket(config RedisConfig) *RedisBucket {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	return &RedisBucket{client: client}
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
	lastRefillCmd := pipe.Get(Ctx, lastRefillKey)
	tokensCmd := pipe.Get(Ctx, tokensKey)
	_, err := pipe.Exec(Ctx)
	if err != nil && err != redis.Nil {
		return
	}

	var lastRefill int
	if lastRefillCmd.Err() == redis.Nil {
		lastRefill = int(now)
		pipe.Set(Ctx, lastRefillKey, now, 0)
	} else {
		lastRefill, err = strconv.Atoi(lastRefillCmd.Val())
		if err != nil {
			return
		}
	}

	var tokens int
	if tokensCmd.Err() == redis.Nil {
		tokens = bucketConfig.Capacity
		pipe.Set(Ctx, tokensKey, bucketConfig.Capacity, 0)
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

	pipe.Set(Ctx, tokensKey, updatedTokens, 0)
	pipe.Set(Ctx, lastRefillKey, now, 0)
	pipe.Exec(Ctx)
}

func (rb *RedisBucket) isAllow(bucketConfig BucketConfig) bool {
	tokens, err := rb.client.Get(Ctx, bucketConfig.Key+":tokens").Int()
	if err != nil && err != redis.Nil {
		return false
	}
	return tokens > 0
}

func (rb *RedisBucket) consumeToken(bucketConfig BucketConfig) bool {
	key := bucketConfig.Key
	tokensKey := key + ":tokens"
	lastRefillKey := key + ":last_refill"

	pipe := rb.client.TxPipeline()
	tokensCmd := pipe.Get(Ctx, tokensKey)
	_, err := pipe.Exec(Ctx)
	if err != nil && err != redis.Nil {
		return false
	}

	tokens, err := strconv.Atoi(tokensCmd.Val())
	if tokens >= 1 {
		pipe.Decr(Ctx, tokensKey)
		pipe.Set(Ctx, lastRefillKey, time.Now().Unix(), 0)
		pipe.Exec(Ctx)
		return true
	}
	return false
}

func (rb *RedisBucket) getRemainingTokens(bucketConfig BucketConfig) int {
	tokens, err := rb.client.Get(Ctx, bucketConfig.Key+":tokens").Int()
	if err != nil && err != redis.Nil {
		return 0
	}
	return tokens
}
