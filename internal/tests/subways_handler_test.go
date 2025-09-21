package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/yeison2020/subway-routing-service/internal/config"
	"github.com/yeison2020/subway-routing-service/internal/mbta"
	"github.com/yeison2020/subway-routing-service/internal/services"
)

type MockClient struct{}

func (m MockClient) FetchRoutes(apiKey string) ([]mbta.Route, error) {
	return []mbta.Route{
		{ID: "Red", Name: "Red Line"},
		{ID: "Green", Name: "Green Line"},
	}, nil
}

func (m MockClient) FetchStops(apiKey, routeID string) ([]mbta.Station, error) {
	if routeID == "Red" {
		return []mbta.Station{{ID: "R1", Name: "South Station"}}, nil
	} else if routeID == "Green" {
		return []mbta.Station{{ID: "G1", Name: "Park Street"}}, nil
	}
	return nil, nil
}
func TestGetSubwaysHandler(t *testing.T) {
	cfg := &config.Config{MBTAApiKey: "fake-key"}

	router := gin.Default()
	router.GET("/subways", services.GetSubwaysHandler(cfg, MockClient{}))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/subways", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	var resp map[string][]mbta.Route
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	routes := resp["routes"]
	if len(routes) != 2 {
		t.Fatalf("expected 2 routes, got %d", len(routes))
	}

	if len(routes[0].Stations) == 0 {
		t.Errorf("expected stations for route %s", routes[0].ID)
	}
}
