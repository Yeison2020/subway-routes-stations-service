package tests


import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yeison2020/subway-routing-service/internal/utils"
)

func TestFindAllRoutesBFS(t *testing.T) {
	// Setup a small graph
	g := &utils.Graph{
		Nodes: map[string][]utils.Edge{
			"A": {{StationID: "B", Line: "Red"}, {StationID: "C", Line: "Blue"}},
			"B": {{StationID: "D", Line: "Red"}},
			"C": {{StationID: "D", Line: "Blue"}},
			"D": {},
		},
		Names: map[string]string{
			"A": "Station A",
			"B": "Station B",
			"C": "Station C",
			"D": "Station D",
		},
	}

	tests := []struct {
		start     string
		end       string
		maxRoutes int
		expected  [][]string
	}{
		{
			start:     "A",
			end:       "D",
			maxRoutes: 2,
			expected: [][]string{
				{"Station A", "Station B", "Station D"},
				{"Station A", "Station C", "Station D"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.start+"-to-"+tt.end, func(t *testing.T) {
			routes, _, err := g.FindAllRoutesBFS(tt.start, tt.end, tt.maxRoutes)
			assert.NoError(t, err)
			assert.ElementsMatch(t, tt.expected, routes)
		})
	}
}

func TestFindAllRoutesBFS_NoRoute(t *testing.T) {
	// Graph with no path between start and end
	g := &utils.Graph{
		Nodes: map[string][]utils.Edge{
			"A": {},
			"B": {},
		},
		Names: map[string]string{
			"A": "Station A",
			"B": "Station B",
		},
	}

	_, _, err := g.FindAllRoutesBFS("A", "B", 1)
	assert.Error(t, err)
}
