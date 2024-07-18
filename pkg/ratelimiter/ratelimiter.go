package ratelimiter

func IsAllow(rlConfig RateLimitConfig) (bool, RateLimitRemaining) {
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
		refillToken(bucketConfig)
		if bucketConfig.Capacity != UNLIMITED_RATE && !isAllow(bucketConfig) {
			allowed = false
		}
	}

	remaining := RateLimitRemaining{
		AccountID:        rlConfig.AccountID,
		ServiceName:      rlConfig.ServiceName,
		RequestPerSecond: getRemainingTokens(bucketConfigs[0]),
		RequestPerMinute: getRemainingTokens(bucketConfigs[1]),
		RequestPerHour:   getRemainingTokens(bucketConfigs[2]),
		RequestPerDay:    getRemainingTokens(bucketConfigs[3]),
		RequestPerWeek:   getRemainingTokens(bucketConfigs[4]),
		RequestPerMonth:  getRemainingTokens(bucketConfigs[5]),
		RequestPerYear:   getRemainingTokens(bucketConfigs[6]),
	}

	return allowed, remaining
}

func UpdateToken(rlConfig RateLimitConfig) {
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
				consumeToken(bucketConfig)
			}
		}
	}()
}
