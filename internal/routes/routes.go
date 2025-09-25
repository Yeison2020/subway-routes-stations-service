package routes

import (
	"net/http"
	"github.com/gin-gonic/gin"

	"github.com/yeison2020/subway-routing-service/internal/config"
	"github.com/yeison2020/subway-routing-service/internal/mbta"
	"github.com/yeison2020/subway-routing-service/internal/services"

	_ "github.com/yeison2020/subway-routing-service/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)
func RegisterRoutes(server *gin.Engine) {

	cfg := config.LoadConfig()

	

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.GET("/api/v1/healthz", services.HealthHandler(cfg, mbta.RealClient{}))

	server.GET("/api/v1/subways", services.GetSubwaysHandler(cfg, mbta.RealClient{}))

	server.GET("/api/v1/routes", services.GetRouteHandler(cfg, mbta.RealClient{}))

	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Endpoint not found",
			"path":    ctx.Request.URL.Path,
			"message": "Please check endpoint",
		})

	})
}
