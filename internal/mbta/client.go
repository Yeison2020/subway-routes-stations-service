package mbta

import (
	"encoding/json"
	"fmt"
	"net/http"

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
    FetchRoutes(apiKey string) ([]Route, error)
    FetchStops(apiKey, routeID string) ([]Station, error)
}

type RealClient struct{}

func (RealClient) FetchRoutes(apiKey string) ([]Route, error) {
    return FetchRoutes(apiKey)
}

func (RealClient) FetchStops(apiKey, routeID string) ([]Station, error) {
    return FetchStopsCached(apiKey, routeID)
}

// FetchRoutes gets all subway routes (type 0=light rail, 1=subway)
func FetchRoutes(apiKey string) ([]Route, error) {
	url := "https://api-v3.mbta.com/routes?filter[type]=0,1"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
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
func FetchStops(apiKey, routeID string) ([]Station, error) {
	url := fmt.Sprintf("https://api-v3.mbta.com/stops?filter[route]=%s", routeID)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
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
func FetchStopsCached(apiKey, routeID string) ([]Station, error) {
	if cached, ok := stationCache.Get(routeID); ok {
		return cached, nil
	}

	stations, err := FetchStops(apiKey, routeID)
	if err != nil {
		return nil, err
	}

	stationCache.Set(routeID, stations)
	return stations, nil
}
