package main

import (


	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/yeison2020/subway-routing-service/internal/config"
	"github.com/yeison2020/subway-routing-service/internal/middlerware"
	"github.com/yeison2020/subway-routing-service/internal/routes"
)

func main() {

    // Load .env
	_ = godotenv.Load()

	// Load config file 

	cfg := config.LoadConfig()

    gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	logger := logrus.New()

	// Plug in middlewares
	r.Use(middlerware.LoggingMiddleware(logger))
	r.Use(middlerware.RequestIDMiddlerware())


	// Routes 
	routes.RegisterRoutes(r)



   	// Run server
	r.Run(":" + cfg.ServerPort)

}
