package main

import (
	"ratelimiter/pkg/ratelimiter"
	"time"

	"github.com/rs/zerolog/log"
)

func main() {
	config := ratelimiter.RateLimitConfig{
		AccountID:        "user123",
		ServiceName:      "A",
		RequestPerSecond: 2,
		RequestPerMinute: 50,
		RequestPerHour:   400,
		RequestPerDay:    ratelimiter.UNLIMITED_RATE,
		RequestPerWeek:   ratelimiter.UNLIMITED_RATE,
		RequestPerMonth:  ratelimiter.UNLIMITED_RATE,
		RequestPerYear:   ratelimiter.UNLIMITED_RATE,
	}

	bucket := ratelimiter.NewDefaultBucket()
	rateLimiter := ratelimiter.NewRateLimiter(bucket)

	for i := 0; i < 1000; i++ {
		allowed, remaining := rateLimiter.IsAllow(config)
		if allowed {
			log.Info().Msgf("Request allowed. Remaining limits: %+v\n", remaining)
			rateLimiter.UpdateToken(config)
		} else {
			log.Info().Msgf("Request denied. Remaining limits: %+v\n", remaining)
		}
		time.Sleep(1 * time.Second)
	}
}
