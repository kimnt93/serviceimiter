# Service Limiter

This Go package implements a rate limiter using the Token Bucket algorithm. It supports Redis and in-memory buckets, as well as manual token management. The rate limiter can be configured by user ID, service name, and various time intervals (second, hour, day, week, month, year), with a default unlimited limit of `-1`.

## Features

- **Token Bucket Algorithm**: Efficient rate limiting.
- **Redis Bucket**: Distributed rate limiting with Redis.
- **In-Memory Bucket**: Single-node or testing use.
- **Manual Token Management**: Manual token updates and checks.

## Installation

Install the package using:

```bash
go get github.com/kimnt93/servicelimiter
```

## Usage

Example usage can be found in the `test/` directory of the repository. The directory contains examples for:

- Using the default in-memory bucket.
- Using the Redis bucket for distributed scenarios.
- Manual token management.

### Quick Example

Here’s a shortened example of using the in-memory bucket:

```go
package main

import (
    "log"
    "time"
    "github.com/kimnt93/servicelimiter"
)

func main() {
    config := ratelimiter.RateLimitConfig{
        AccountID:        "user123",
        ServiceName:      ratelimiter.AUTO_SERVICE_NAME,
        RequestPerSecond: 2,
        RequestPerMinute: 27,
        RequestPerHour:   400,
        RequestPerDay:    ratelimiter.UNLIMITED_RATE,
        RequestPerWeek:   ratelimiter.UNLIMITED_RATE,
        RequestPerMonth:  ratelimiter.UNLIMITED_RATE,
        RequestPerYear:   ratelimiter.UNLIMITED_RATE,
    }

    bucket := ratelimiter.NewDefaultBucket()
    rateLimiter := ratelimiter.NewRateLimiter(bucket)

    for i := 0; i < 1000; i++ {
        allowed, remaining, _, err := rateLimiter.Run(config, PrintOnly)
        log.Printf("Request allowed: %t, Remaining limits: %+v, Error: %+v\n", allowed, remaining, err)
        time.Sleep(1 * time.Second)
    }
}
```

## Configuration

- **AccountID**: Unique identifier for the user or account.
- **ServiceName**: Service name for the rate limit.
- **RequestPerSecond**: Requests per second.
- **RequestPerMinute**: Requests per minute.
- **RequestPerHour**: Requests per hour.
- **RequestPerDay**: Requests per day.
- **RequestPerWeek**: Requests per week.
- **RequestPerMonth**: Requests per month.
- **RequestPerYear**: Requests per year.
- **UNLIMITED_RATE**: Set to `-1` for unlimited rate.

## License

Licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
