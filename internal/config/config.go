package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerPort string
	MBTAApiKey string
	Service    string
	Env        string
	Version    string
}

func LoadConfig() *Config {

	port := os.Getenv("PORT")

	/*
	Todo:check missing values and add defauls
	*/

	serviceName := os.Getenv("DD_SERVICE")

	envName := os.Getenv("DD_ENV")

	version := os.Getenv("DD_VERSION")

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
		Service:    serviceName,
		Env:        envName,
		Version: version,
	}
}
