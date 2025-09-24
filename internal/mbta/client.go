package mbta

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	httptrace "github.com/DataDog/dd-trace-go/contrib/net/http/v2"
)

// Models

type Station struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Route struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Stations []Station `json:"stations"`
}

// MBTA API responses
type routesResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Attributes struct {
			Name string `json:"long_name"`
		} `json:"attributes"`
	} `json:"data"`
}

type stopsResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Attributes struct {
			Name string `json:"name"`
		} `json:"attributes"`
	} `json:"data"`
}

// Testing:

type Client interface {
	FetchRoutes(apiKey string, ctx *gin.Context) ([]Route, error)
	FetchStops(apiKey, routeID string, ctx *gin.Context) ([]Station, error)
}

type RealClient struct{}

func (RealClient) FetchRoutes(apiKey string, ctx *gin.Context) ([]Route, error) {
	return FetchRoutes(apiKey, ctx)
}

func (RealClient) FetchStops(apiKey, routeID string, ctx *gin.Context) ([]Station, error) {
	return FetchStopsCached(apiKey, routeID, ctx)
}

// FetchRoutes gets all subway routes (type 0=light rail, 1=subway)
func FetchRoutes(apiKey string, ctx *gin.Context) ([]Route, error) {
	url := "https://api-v3.mbta.com/routes?filter[type]=0,1"

	req, _ := http.NewRequestWithContext(ctx.Request.Context(), "GET", url, nil)

	req.Header.Set("x-api-key", apiKey)

	client := httptrace.WrapClient(http.DefaultClient)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r routesResponse

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	routes := []Route{}
	for _, d := range r.Data {
		routes = append(routes, Route{
			ID:   d.ID,
			Name: d.Attributes.Name,
		})
	}
	return routes, nil
}

// FetchStops gets all stations for a route from MBTA API
func FetchStops(apiKey, routeID string, ctx *gin.Context) ([]Station, error) {
	url := fmt.Sprintf("https://api-v3.mbta.com/stops?filter[route]=%s", routeID)

	req, _ := http.NewRequestWithContext(ctx.Request.Context(), "GET", url, nil)

	fmt.Print(ctx)
	req.Header.Set("x-api-key", apiKey)
	client := httptrace.WrapClient(http.DefaultClient)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r stopsResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	stations := []Station{}
	for _, d := range r.Data {
		stations = append(stations, Station{
			ID:   d.ID,
			Name: d.Attributes.Name,
		})
	}
	return stations, nil
}

// FetchStopsCached returns stations from cache or MBTA if not cached
func FetchStopsCached(apiKey, routeID string, ctx *gin.Context) ([]Station, error) {
	if cached, ok := stationCache.Get(routeID); ok {
		return cached, nil
	}

	stations, err := FetchStops(apiKey, routeID, ctx)
	if err != nil {
		return nil, err
	}

	stationCache.Set(routeID, stations)
	return stations, nil
}
