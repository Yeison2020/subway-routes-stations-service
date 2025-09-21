package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/yeison2020/subway-routing-service/internal/config"
	"github.com/yeison2020/subway-routing-service/internal/mbta"
	"github.com/yeison2020/subway-routing-service/internal/services"
)


func RegisterRoutes(server *gin.Engine) {

	cfg := config.LoadConfig()

	server.GET("/api/v1/health", services.HealthHandler(cfg, mbta.RealClient{}))

	server.GET("/api/v1/subways", services.GetSubwaysHandler(cfg, mbta.RealClient{}))

    server.GET("/api/v1/routes", services.GetRouteHandler(cfg, mbta.RealClient{}))

	server.Run("localhost:8080")

	fmt.Println("Running in localhost port 8080")
}
