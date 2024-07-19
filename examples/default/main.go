package main

import (
	"errors"
	"ratelimiter/pkg/servicelimiter"
	"time"

	"github.com/rs/zerolog/log"
)

func PrintOnly() (int, int, int) {
	return 1, 2, 3
}

func PrintArg(message string, number int) (string, int, error) {
	return message, number, errors.New("this is default error")
}

func main() {
	config1 := servicelimiter.RateLimitConfig{
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

	config2 := servicelimiter.RateLimitConfig{
		AccountID:        "user123",
		ServiceName:      servicelimiter.AUTO_SERVICE_NAME,
		RequestPerSecond: 2,
		RequestPerMinute: 20,
		RequestPerHour:   400,
		RequestPerDay:    servicelimiter.UNLIMITED_RATE,
		RequestPerWeek:   servicelimiter.UNLIMITED_RATE,
		RequestPerMonth:  servicelimiter.UNLIMITED_RATE,
		RequestPerYear:   servicelimiter.UNLIMITED_RATE,
	}

	bucket := servicelimiter.NewDefaultBucket()
	rl := servicelimiter.NewRateLimiter(bucket)

	for i := 0; i < 1000; i++ {
		allowed, remaining, funcRt, err := rl.Run(config1, PrintOnly)
		log.Info().Msgf("Request allowed: %t, Remaining limits: %+v, Function return: %+v, Error: %+v\n", allowed, remaining, funcRt, err)

		allowed, remaining, funcRt, err = rl.Run(config2, PrintArg, "Hello worldddd", 10)
		log.Info().Msgf("Request allowed: %t, Remaining limits: %+v, Function return: %+v, Error: %+v\n", allowed, remaining, funcRt, err)

		time.Sleep(1 * time.Second)
	}
}
