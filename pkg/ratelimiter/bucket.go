package ratelimiter

type Bucket interface {
	refillToken(bucketConfig BucketConfig)
	isAllow(bucketConfig BucketConfig) bool
	consumeToken(bucketConfig BucketConfig) bool
	getRemainingTokens(bucketConfig BucketConfig) int
}
