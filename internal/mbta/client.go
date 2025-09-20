package mbta

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// MBTA API Functions

func FetchRoutes(apikey string) ([]Route, error){

	url := "https://api-v3.mbta.com/routes?filter[type]=0,1"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("x-api-key", apikey)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return nil , err
	}

	defer resp.Body.Close()

	var r RouteResponse

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	
	}

	routes := []Route{}

	for _, d := range r.Data {
		routes = append(routes, Route{
			ID: d.ID,
			Name: d.Attributes.Name,
		})
	}

	return routes, nil

}


// FetchStops gets all stations for a route from MBTA API

func FetchStops(apikey, routeID string) ([]Station, error){

	url := fmt.Sprintf("https://api-v3.mbta.com/stops?filter[route]=%s", routeID)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("x-api-key", apikey)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var r stopsResponse

	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	stations := []Station{}

	for _, d := range r.Data{
		stations = append(stations, Station{
			ID: d.ID,
			Name: d.Attributes.Name,
		})
	}

 return  stations, nil

}


func FetchStopsCached(apikey, routeID string) ([]Station, error){
	
	if cached, ok := StationCache.Get(routeID); ok {
		return cached, nil
	}

	stations, err := FetchStops(apikey, routeID)
	if err != nil {
		return nil, err
	}
	StationCache.Set(routeID, stations)
	return stations, nil
}