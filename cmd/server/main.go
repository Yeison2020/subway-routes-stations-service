package main

import (
	
     "log"
	 "net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/yeison2020/subway-routing-service/internal/config"
	"github.com/yeison2020/subway-routing-service/internal/routes"
	"github.com/yeison2020/subway-routing-service/internal/mbta"
	"github.com/yeison2020/subway-routing-service/internal/middlerware"


)

func main() {

   // Load .env
	_ = godotenv.Load()

	// Initialize Zap logger
	middlerware.InitLogger() // returns *zap.Logger

	// Init MBTA local cache
	mbta.InitCache(300) // TTL 5 minutes

	cfg := config.LoadConfig()

    gin.SetMode(gin.ReleaseMode)

	server := gin.New()

	// Plug middlerware
	server.Use(middlerware.LoggerMiddleware(middlerware.Logger))
	server.Use(middlerware.RequestIdMiddleware())
	server.Use(middlerware.PanicHandler())


	// Routes 
	routes.RegisterRoutes(server)

	// 404 routes
	server.NoRoute(func(ctx *gin.Context){
    	ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Endpoint not found",
			"path":    ctx.Request.URL.Path,
			"message": "Please check endpoint",
		})

	})

		// App startup logs
	middlerware.Logger.Info("Application started")
	middlerware.Logger.Info("Server running",
		zap.String("port", cfg.ServerPort),
		zap.String("mode", "release"),
	)



   	// Run server
	log.Fatal(server.Run(":" + cfg.ServerPort))

}
