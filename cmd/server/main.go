package main

import (
     "log"
	 "fmt"
	 "net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/yeison2020/subway-routing-service/internal/config"
	"github.com/yeison2020/subway-routing-service/internal/routes"
	"github.com/yeison2020/subway-routing-service/internal/middlerware"
	"github.com/yeison2020/subway-routing-service/internal/mbta"
)

func main() {

    // Load .env
	_ = godotenv.Load()
	// Logger initialization
	middlerware.InitLogger()

	// Init MBTA local cache
	mbta.InitCache(300) // TTL 5 minutes

	// Load config file 
	cfg := config.LoadConfig()

    gin.SetMode(gin.ReleaseMode)

	server := gin.New()

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

	fmt.Println("Server started running")


   	// Run server
	log.Fatal(server.Run(":" + cfg.ServerPort))

}
