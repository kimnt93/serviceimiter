package main

import (
	"servicelimiter/pkg/servicelimiter"
	"time"

	"github.com/rs/zerolog/log"
)

func PrintOnly() (int, int, int) {
	return 1, 2, 3
}

func main() {
	config := servicelimiter.RateLimitConfig{
		AccountID:        "user123",
		ServiceName:      servicelimiter.AUTO_SERVICE_NAME,
		RequestPerSecond: 2,
		RequestPerMinute: 27,
		RequestPerHour:   400,
		RequestPerDay:    servicelimiter.UNLIMITED_RATE,
		RequestPerWeek:   servicelimiter.UNLIMITED_RATE,
		RequestPerMonth:  servicelimiter.UNLIMITED_RATE,
		RequestPerYear:   servicelimiter.UNLIMITED_RATE,
	}

	bucket := servicelimiter.NewDefaultBucket() // or can use servicelimiter.NewRedisBucket("localhost:6379", "", 0)
	rl := servicelimiter.NewRateLimiter(bucket)

	for i := 0; i < 1000; i++ {
		allowed, remaining := rl.IsAllow(config)
		if allowed {
			log.Info().Msgf("Request allowed: %t, Remaining limits: %+v", allowed, remaining)
			// Start function you want
			PrintOnly()
			// Update token
			rl.UpdateToken(config)
		} else {
			log.Info().Msgf("RL exceeded: %t, Remaining limits: %+v", allowed, remaining)

		}

		log.Info().Msgf("Request allowed: %t, Remaining limits: %+v", allowed, remaining)

		time.Sleep(1 * time.Second)
	}
}
