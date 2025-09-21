package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yeison2020/subway-routing-service/internal/config"
	"github.com/yeison2020/subway-routing-service/internal/mbta"
	"github.com/yeison2020/subway-routing-service/internal/middlerware"
	"github.com/yeison2020/subway-routing-service/internal/utils"
)


func GetRouteHandler(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := ctx.Query("start")
		end := ctx.Query("end")

		if start == "" || end == "" {
			middlerware.Logger.Error("Empty value - start and end parameters required")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "start and end parameters required"})
			return
		}

		// Fetch all routes
		routes, err := mbta.FetchRoutes(cfg.MBTAApiKey)

		if err != nil {
			middlerware.Logger.Error(err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}


		// Fetch stations for each route (cached)

		for i, r := range routes {
			stations, err := mbta.FetchStopsCached(cfg.MBTAApiKey, r.ID)

			if err != nil {
				middlerware.Logger.Error(err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			routes[i].Stations = stations
		}


		// Build the graph
		g := utils.BuildGraph(routes)

		// Find the route

		stationsResult, linesResult, err := g.FindRouteBFS(start, end)


		if err != nil {
			middlerware.Logger.Error(err.Error())
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return 
		}

		// BuildRouteDescription here

		description := utils.BuildRouteDescription(stationsResult, linesResult)

		//Return JSON response

		ctx.JSON(http.StatusOK, gin.H{
			"stations":    stationsResult,
			"lines":       linesResult,
			"description": description, // Human-readable route
		})
	}


}