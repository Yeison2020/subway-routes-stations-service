package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yeison2020/subway-routing-service/internal/config"
	"github.com/yeison2020/subway-routing-service/internal/mbta"
)

// GetSubwaysHandlerSwagger godoc
// @Summary Get all subway routes with their stations
// @Description Returns a list of subway routes with nested stations
// @Tags Subways
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /subways [get]
func GetSubwaysHandlerSwagger(c *gin.Context) {
	c.JSON(200, map[string]interface{}{
		"routes": []map[string]interface{}{
			{
				"id":   "Red",
				"name": "Red Line",
				"stations": []map[string]string{
					{"id": "place-alfcl", "name": "Alewife"},
					{"id": "place-davis", "name": "Davis"},
				},
			},
			{
				"id":   "Green",
				"name": "Green Line",
				"stations": []map[string]string{
					{"id": "place-lech", "name": "Lechmere"},
					{"id": "place-north", "name": "North Station"},
				},
			},
		},
	})
}

// GetSubwaysHandler returns all subway routes with their stations
func GetSubwaysHandler(cfg *config.Config, client mbta.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. Fetch all subway routes
		routesData, err := client.FetchRoutes(cfg.MBTAApiKey, ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error 1": err.Error()})
			return
		}
		var routes []mbta.Route

		// 2. Fetch stations for each route and nest
		for _, r := range routesData {
			stations, err := client.FetchStops(cfg.MBTAApiKey, r.ID, ctx)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error 2": err.Error()})
				return
			}

			routes = append(routes, mbta.Route{
				ID:       r.ID,
				Name:     r.Name,
				Stations: stations,
			})
		}

		// 3. Return nested response
		ctx.JSON(http.StatusOK, gin.H{
			"routes": routes,
		})
	}
}
