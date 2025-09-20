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
)

func main() {

    // Load .env
	_ = godotenv.Load()
	// Load config file 
	cfg := config.LoadConfig()
	// Logger initialization
	middlerware.InitLogger()

    gin.SetMode(gin.ReleaseMode)

	server := gin.New()
	server.Use(middlerware.LoggerMiddleware(middlerware.Logger))

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
