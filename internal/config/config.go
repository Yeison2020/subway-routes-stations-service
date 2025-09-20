package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerPort string
	MBTAApiKey string
}


func LoadConfig() *Config {

	port := os.Getenv("PORT")

	if port == "" {

		fmt.Println("Deaulting to port 8080")

		port = "8080"

	}

	mbtaKey := os.Getenv("MBTA_API_KEY")

	if mbtaKey == "" {
		fmt.Println("Missing MBTA KEY")
	}

	return &Config{
		ServerPort: port,
		MBTAApiKey: mbtaKey,
	}
}
