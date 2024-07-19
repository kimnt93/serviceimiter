package servicelimiter

type BucketConfig struct {
	Key        string
	PeriodType int
	Capacity   int
	Ttl        int
}

type Bucket interface {
	refillToken(bucketConfig BucketConfig)
	isAllow(bucketConfig BucketConfig) bool
	consumeToken(bucketConfig BucketConfig) bool
	getRemainingTokens(bucketConfig BucketConfig) int
}

var isBucketInitialized = false
