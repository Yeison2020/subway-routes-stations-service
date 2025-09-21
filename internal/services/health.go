package services

import (
   "fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yeison2020/subway-routing-service/internal/config"
	"github.com/yeison2020/subway-routing-service/internal/mbta"
)




func HealthHandler(cfg *config.Config, client mbta.Client)  gin.HandlerFunc{

	return func(ctx *gin.Context) {

		health := map[string]interface{} {
			"status": "ok",
		}

		routes, err := client.FetchRoutes(cfg.MBTAApiKey)

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