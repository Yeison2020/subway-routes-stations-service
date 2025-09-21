package routes

import (

	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/yeison2020/subway-routing-service/internal/services"
	"github.com/yeison2020/subway-routing-service/internal/config"

)


func RegisterRoutes(server *gin.Engine) {

	cfg := config.LoadConfig()

	server.GET("/api/v1/health", services.HealthHandler(cfg))

	server.GET("/api/v1/subways", services.GetSubwaysHandler(cfg))

    server.GET("/api/v1/routes", services.GetRouteHandler(cfg))

	server.Run("localhost:8080")

	fmt.Println("Running in localhost port 8080")
}
