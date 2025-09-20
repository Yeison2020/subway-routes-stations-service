package services

import (
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yeison2020/subway-routing-service/internal/config"
	"github.com/yeison2020/subway-routing-service/internal/mbta"
)

func GetSubwaysHandler(cfg *config.Config) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		routesData, err := mbta.FetchRoutes(cfg.MBTAApiKey)

		fmt.Print(routesData)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error 1": err.Error()})
			return 
		}

		var routes []mbta.Route

		for _, r := range routesData {
			stations, err := mbta.FetchStopsCached(cfg.MBTAApiKey, r.ID)

			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error 1": err.Error()})
				return 
			}
			routes = append(routes, mbta.Route{
				ID: r.ID,
				Name: r.Name,
				Stations: stations,
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"routes": routes,
		})

	}

}
