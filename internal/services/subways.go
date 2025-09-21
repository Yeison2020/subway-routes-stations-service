package services

import (
	"net/http"


	"github.com/gin-gonic/gin"
	"github.com/yeison2020/subway-routing-service/internal/config"
	"github.com/yeison2020/subway-routing-service/internal/mbta"
)



// GetSubwaysHandler returns all subway routes with their stations
func GetSubwaysHandler(cfg *config.Config,  client mbta.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Fetch all subway routes
		routesData, err := client.FetchRoutes(cfg.MBTAApiKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error 1": err.Error()})
			return
		}

		var routes []mbta.Route

		// 2. Fetch stations for each route and nest
		for _, r := range routesData {
			stations, err := client.FetchStops(cfg.MBTAApiKey, r.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error 2": err.Error()})
				return
			}

			routes = append(routes, mbta.Route{
				ID:       r.ID,
				Name:     r.Name,
				Stations: stations,
			})
		}

		// 3. Return nested response
		c.JSON(http.StatusOK, gin.H{
			"routes": routes,
		})
	}
}
