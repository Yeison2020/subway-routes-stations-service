package config

import (
	"log"
	"os"

	"github.com/yeison2020/subway-routing-service/internal/helper"
	
)

type Config struct {
	ServerPort string
	MBTAApiKey string
	Service    string
	Env        string
	Version    string
}

func LoadConfig() *Config {

	port := helper.GetEnv("PORT", "8080")
	serviceName := helper.GetEnv("DD_SERVICE","default-service-name" )
	envName :=  helper.GetEnv("DD_ENV",  "dev")
	version :=  helper.GetEnv("DD_VERSION", "1.0")

	mbtaKey := os.Getenv("MBTA_API_KEY")
	if mbtaKey == "" {
		log.Fatal("Missing MBTA KEY")
	}

	return &Config{
		ServerPort: port,
		MBTAApiKey: mbtaKey,
		Service:    serviceName,
		Env:        envName,
		Version: version,
	}
}
