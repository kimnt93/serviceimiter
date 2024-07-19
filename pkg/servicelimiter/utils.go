package servicelimiter

import "fmt"

func getKey(accountID string, serviceName string, period int) string {
	return fmt.Sprintf("%s:%s:%s:%d", RATE_LIMIT_PREFIX, accountID, serviceName, period)
}
