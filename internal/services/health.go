package services

import (


	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yeison2020/subway-routing-service/internal/config"
)


func HealthHandler(cfg *config.Config)  gin.HandlerFunc{

	return func(ctx *gin.Context) {

		ctx.JSON(http.StatusOK, gin.H{"status ok - Todo": "need work", "message":"Working"})
		// health := map[string]interface{} {
		// 	"status": "ok",
		// }

		// routes, err := mbta.FetchRoutes(cfg.MBTAApiKey)


		// if err != nil {
		// 	ctx.JSON(http.StatusBadRequest, gin.H(err.Error()))
		// 	return
		// } else {
		// 	health["mbta"] = fmt.Sprinf("ok, %d routes", len(routes))
		// }


		// code := http.StatusOK 
		// if health["status"] == "fail" {
		// 	code = http.StatusServiceUnavailable
		// }

		// ctx.JSON(code, health)

	}

}