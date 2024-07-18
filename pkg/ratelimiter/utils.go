package ratelimiter

import "fmt"

func GetKey(accountID string, serviceName string, period int) string {
	return fmt.Sprintf("rate_limit:%s:%s:%d", accountID, serviceName, period)
}
