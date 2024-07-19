package servicelimiter

type RateLimitConfig struct {
	AccountID        string
	ServiceName      string
	RequestPerSecond int
	RequestPerMinute int
	RequestPerHour   int
	RequestPerDay    int
	RequestPerWeek   int
	RequestPerMonth  int
	RequestPerYear   int
}

type RateLimitRemaining struct {
	AccountID        string
	ServiceName      string
	RequestPerSecond int
	RequestPerMinute int
	RequestPerHour   int
	RequestPerDay    int
	RequestPerWeek   int
	RequestPerMonth  int
	RequestPerYear   int
}
