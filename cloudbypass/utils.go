package cloudbypass

import "os"

func getEnv(key string, default_value string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return default_value
}
