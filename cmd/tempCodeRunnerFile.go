package main

import (
	"fmt"
	"ratelimiter/pkg/ratelimiter"
	"time"
)

func main() {
	config := ratelimiter.RateLimitConfig{
		AccountID:        "user123",
		ServiceName:      "A",
		RequestPerSecond: 2,
		RequestPerMinute: 10,
		RequestPerHour:   600,
		RequestPerDay:    ratelimiter.UNLIMITED_RATE,
		RequestPerWeek:   ratelimiter.UNLIMITED_RATE,
		RequestPerMonth:  ratelimiter.UNLIMITED_RATE,
		RequestPerYear:   ratelimiter.UNLIMITED_RATE,
	}

	ratelimiter.NewDefaultTokenBucket()

	for i := 0; i < 1000; i++ {
		allowed, remaining := ratelimiter.IsAllow(config)
		if allowed {
			fmt.Printf("Request allowed. Remaining limits: %+v\n", remaining)
			ratelimiter.UpdateToken(config)
		} else {
			fmt.Printf("Request denied. Remaining limits: %+v\n", remaining)
		}
		time.Sleep(1 * time.Second)
	}
}
