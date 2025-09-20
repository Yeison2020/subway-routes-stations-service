package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"

	// "github.com/yeison2020/subway-routing-service/internal/middlerware"
	"github.com/yeison2020/subway-routing-service/internal/services"
	"github.com/yeison2020/subway-routing-service/internal/config"

)

func RegisterRoutes(server *gin.Engine) {


	server.GET("/api/v1/health", services.HealthHandler(&config.Config{}))

	server.GET("/api/v1/subways", services.GetSubwaysHandler(&config.Config{}))


	// autheticated := server.Group("/")

	// autheticated.Use(middleware.Autheticate)
	
	// autheticated.POST("/events", CreateEvents)
	// autheticated.PUT("/events/:id",  UpdateEvent)
	// autheticated.DELETE("/events/:id", DeleteEvent)
	// autheticated.POST("/events/:id/register", registerForEvent)
	// autheticated.DELETE("/events/:id/register", cancelERegistration)

	// server.POST("/signup", signup)
	// server.POST("/login", login)

	server.Run("localhost:8080")

	fmt.Println("Running in localhost port 8080")
}
