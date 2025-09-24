package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/yeison2020/subway-routing-service/internal/config"
	"github.com/yeison2020/subway-routing-service/internal/mbta"
	"github.com/yeison2020/subway-routing-service/internal/services"

	_ "github.com/yeison2020/subway-routing-service/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Example API
// @version 1.0
// @description This is a sample server.
// @host localhost:8080
// @BasePath /api/v1

func RegisterRoutes(server *gin.Engine) {

	cfg := config.LoadConfig()

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.GET("/api/v1/health", services.HealthHandler(cfg, mbta.RealClient{}))

	server.GET("/api/v1/subways", services.GetSubwaysHandler(cfg, mbta.RealClient{}))

	server.GET("/api/v1/routes", services.GetRouteHandler(cfg, mbta.RealClient{}))

	fmt.Println("Running in localhost port 8080")
}
