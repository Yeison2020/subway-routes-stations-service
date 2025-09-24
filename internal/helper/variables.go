package helper

import (
	"os"
	"log"
)

func GetEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	log.Println("Empty value for", key)
	return defaultVal
}