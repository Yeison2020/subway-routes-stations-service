package services

import (
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yeison2020/subway-routing-service/internal/config"
	"github.com/yeison2020/subway-routing-service/internal/mbta"
	"github.com/yeison2020/subway-routing-service/internal/middlerware"
	"github.com/yeison2020/subway-routing-service/internal/utils"
)

type RouteResponse struct {
	Stations    []string `json:"stations"`
	Lines       []string `json:"lines"`
	Description string   `json:"description"`
}

// GetRouteHandlerSwagger godoc
// @Summary Get subway route
// @Description Returns stations, lines, and description for a route between start and end
// @Tags Routes
// @Accept  json
// @Produce  json
// @Param start query string true "Start station ID"
// @Param end query string true "End station ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /routes [get]
func GetRouteHandlerSwagger(ctx *gin.Context) {
	ctx.JSON(200, map[string]interface{}{
		"stations":    [][]string{{"station1", "station2"}},
		"lines":       [][]string{{"Red Line"}},
		"description": "Take Red Line from station1 to station2",
	})
}

func GetRouteHandler(cfg *config.Config, client mbta.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := ctx.Query("start")
		end := ctx.Query("end")

		if start == "" || end == "" {
			middlerware.Logger.Error("Empty value - start and end parameters required")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "start and end parameters required"})
			return
		}

		// Fetch all routes

		routes, err := client.FetchRoutes(cfg.MBTAApiKey, ctx)

		if err != nil {
			middlerware.Logger.Error(err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Fetch stations for each route (cached)

		for i, r := range routes {
			stations, err := client.FetchStops(cfg.MBTAApiKey, r.ID, ctx)

			if err != nil {
				middlerware.Logger.Error(err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			routes[i].Stations = stations
		}


		fmt.Print(routes)
		// Build the graph
		g := utils.BuildGraph(routes)

		fmt.Println(g)

		// Find the route

		stationsResult, linesResult, err := g.FindAllRoutesBFS(start, end, 3)

		if err != nil {
			middlerware.Logger.Error(err.Error())
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		// BuildRouteDescription here

		description := utils.BuildRouteDescriptions(stationsResult, linesResult)

		//Return JSON response

		ctx.JSON(http.StatusOK, gin.H{
			"stations":    stationsResult,
			"lines":       linesResult,
			"description": description, // Human-readable route
		})
	}

}
