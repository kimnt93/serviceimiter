package ratelimiter

type RateLimiter struct {
	bucket Bucket
}

func NewRateLimiter(bucket Bucket) *RateLimiter {
	return &RateLimiter{bucket: bucket}
}

func (rl *RateLimiter) IsAllow(rlConfig RateLimitConfig) (bool, RateLimitRemaining) {
	allowed := true

	secondKey := GetKey(rlConfig.AccountID, rlConfig.ServiceName, SECOND_PERIOD)
	minuteKey := GetKey(rlConfig.AccountID, rlConfig.ServiceName, MINUTE_PERIOD)
	hourKey := GetKey(rlConfig.AccountID, rlConfig.ServiceName, HOUR_PERIOD)
	dayKey := GetKey(rlConfig.AccountID, rlConfig.ServiceName, DAY_PERIOD)
	weekKey := GetKey(rlConfig.AccountID, rlConfig.ServiceName, WEEK_PERIOD)
	monthKey := GetKey(rlConfig.AccountID, rlConfig.ServiceName, MONTH_PERIOD)
	yearKey := GetKey(rlConfig.AccountID, rlConfig.ServiceName, YEAR_PERIOD)

	bucketConfigs := []BucketConfig{
		{secondKey, SECOND_PERIOD, rlConfig.RequestPerSecond, SECOND_TO_SECOND},
		{minuteKey, MINUTE_PERIOD, rlConfig.RequestPerMinute, MINUTE_TO_SECOND},
		{hourKey, HOUR_PERIOD, rlConfig.RequestPerHour, HOUR_TO_SECOND},
		{dayKey, DAY_PERIOD, rlConfig.RequestPerDay, DAY_TO_SECOND},
		{weekKey, WEEK_PERIOD, rlConfig.RequestPerWeek, WEEK_TO_SECOND},
		{monthKey, MONTH_PERIOD, rlConfig.RequestPerMonth, MONTH_TO_SECOND},
		{yearKey, YEAR_PERIOD, rlConfig.RequestPerYear, YEAR_TO_SECOND},
	}

	for _, bucketConfig := range bucketConfigs {
		// Refill bucket
		rl.bucket.refillToken(bucketConfig)
		if bucketConfig.Capacity != UNLIMITED_RATE && !rl.bucket.isAllow(bucketConfig) {
			allowed = false
		}
	}

	remaining := RateLimitRemaining{
		AccountID:        rlConfig.AccountID,
		ServiceName:      rlConfig.ServiceName,
		RequestPerSecond: rl.bucket.getRemainingTokens(bucketConfigs[0]),
		RequestPerMinute: rl.bucket.getRemainingTokens(bucketConfigs[1]),
		RequestPerHour:   rl.bucket.getRemainingTokens(bucketConfigs[2]),
		RequestPerDay:    rl.bucket.getRemainingTokens(bucketConfigs[3]),
		RequestPerWeek:   rl.bucket.getRemainingTokens(bucketConfigs[4]),
		RequestPerMonth:  rl.bucket.getRemainingTokens(bucketConfigs[5]),
		RequestPerYear:   rl.bucket.getRemainingTokens(bucketConfigs[6]),
	}

	return allowed, remaining
}

func (rl *RateLimiter) UpdateToken(rlConfig RateLimitConfig) {
	go func() {
		secondKey := GetKey(rlConfig.AccountID, rlConfig.ServiceName, SECOND_PERIOD)
		minuteKey := GetKey(rlConfig.AccountID, rlConfig.ServiceName, MINUTE_PERIOD)
		hourKey := GetKey(rlConfig.AccountID, rlConfig.ServiceName, HOUR_PERIOD)
		dayKey := GetKey(rlConfig.AccountID, rlConfig.ServiceName, DAY_PERIOD)
		weekKey := GetKey(rlConfig.AccountID, rlConfig.ServiceName, WEEK_PERIOD)
		monthKey := GetKey(rlConfig.AccountID, rlConfig.ServiceName, MONTH_PERIOD)
		yearKey := GetKey(rlConfig.AccountID, rlConfig.ServiceName, YEAR_PERIOD)

		bucketConfigs := []BucketConfig{
			{secondKey, SECOND_PERIOD, rlConfig.RequestPerSecond, SECOND_TO_SECOND},
			{minuteKey, MINUTE_PERIOD, rlConfig.RequestPerMinute, MINUTE_TO_SECOND},
			{hourKey, HOUR_PERIOD, rlConfig.RequestPerHour, HOUR_TO_SECOND},
			{dayKey, DAY_PERIOD, rlConfig.RequestPerDay, DAY_TO_SECOND},
			{weekKey, WEEK_PERIOD, rlConfig.RequestPerWeek, WEEK_TO_SECOND},
			{monthKey, MONTH_PERIOD, rlConfig.RequestPerMonth, MONTH_TO_SECOND},
			{yearKey, YEAR_PERIOD, rlConfig.RequestPerYear, YEAR_TO_SECOND},
		}

		for _, bucketConfig := range bucketConfigs {
			if bucketConfig.Capacity > 0 {
				rl.bucket.consumeToken(bucketConfig)
			}
		}
	}()
}
