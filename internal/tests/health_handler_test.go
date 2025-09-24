package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/yeison2020/subway-routing-service/internal/config"
	"github.com/yeison2020/subway-routing-service/internal/mbta"
	"github.com/yeison2020/subway-routing-service/internal/services"
)

type mockMBTAClient struct {
	routes []mbta.Route
	err    error
}

func (m *mockMBTAClient) FetchRoutes(apiKey string) ([]mbta.Route, error) {
	return m.routes, m.err
}

func (m *mockMBTAClient) FetchStops(apiKey, routeID string) ([]mbta.Station, error) {
	return []mbta.Station{
		{ID: "1", Name: "Station A"},
		{ID: "2", Name: "Station B"},
	}, nil
}

func TestHealthHandler_Success(t *testing.T) {
	// Arrange
	mockClient := &mockMBTAClient{
		routes: []mbta.Route{{ID: "Red"}, {ID: "Green"}},
		err:    nil,
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	cfg := &config.Config{MBTAApiKey: "fake-key"}

	r.GET("/health", services.HealthHandler(cfg, mockClient))

	// Act
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, "ok", body["status"])
	assert.Contains(t, body["mbta"], "2 routes")
}

func TestHealthHandler_Failure(t *testing.T) {
	mockClient := &mockMBTAClient{
		routes: nil,
		err:    assert.AnError,
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	cfg := &config.Config{MBTAApiKey: "fake-key"}

	// Only register once
	r.GET("/health", services.HealthHandler(cfg, mockClient))

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var body map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Contains(t, body["error"], "assert.AnError")
}
