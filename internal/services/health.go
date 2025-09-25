package services

import (
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yeison2020/subway-routing-service/internal/config"
	"github.com/yeison2020/subway-routing-service/internal/mbta"
)

// HealthHandler godoc
// @Summary Health check
// @Description Returns the status of the service and MBTA API
// @Tags Health
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 503 {object} map[string]interface{}
// @Router /healthz [get]
func HealthHandlerSwagger(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"mbta": "ok, n routes",
	})
}

func HealthHandler(cfg *config.Config, client mbta.Client) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		health := map[string]interface{}{
			"status":  "healthy",
			"message": "API is running",
		}

		routes, err := client.FetchRoutes(cfg.MBTAApiKey, ctx)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			health["mbta"] = fmt.Sprintf("ok, %d routes", len(routes))
		}

		code := http.StatusOK
		if health["status"] == "fail" {
			code = http.StatusServiceUnavailable
		}

		ctx.JSON(code, health)

	}

}
