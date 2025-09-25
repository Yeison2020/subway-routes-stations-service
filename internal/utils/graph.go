package utils

import (
	"fmt"
	"github.com/yeison2020/subway-routing-service/internal/mbta"
)

type Edge struct {
	StationID string
	Line      string
}

type Graph struct {
	Nodes map[string][]Edge
	Names map[string]string
}

// BuildGraph builds a bidirectional graph and adds transfer edges at multi-line stations.
func BuildGraph(routes []mbta.Route) *Graph {
	g := &Graph{
		Nodes: make(map[string][]Edge),
		Names: make(map[string]string),
	}

	// Keep track of which lines each station belongs to
	stationLines := make(map[string][]string)

	for _, route := range routes {
		for i, s := range route.Stations {
			g.Names[s.ID] = s.Name
			stationLines[s.ID] = append(stationLines[s.ID], route.Name)

			// Connect consecutive stations in both directions
			if i < len(route.Stations)-1 {
				next := route.Stations[i+1]
				g.Nodes[s.ID] = append(g.Nodes[s.ID], Edge{StationID: next.ID, Line: route.Name})
				g.Nodes[next.ID] = append(g.Nodes[next.ID], Edge{StationID: s.ID, Line: route.Name})
			}
		}
	}

	// Add transfer edges: allow switching lines at the same station
	for stationID, lines := range stationLines {
		if len(lines) <= 1 {
			continue
		}
		for i := 0; i < len(lines); i++ {
			for j := i + 1; j < len(lines); j++ {
				g.Nodes[stationID] = append(g.Nodes[stationID], Edge{StationID: stationID, Line: lines[j]})
				g.Nodes[stationID] = append(g.Nodes[stationID], Edge{StationID: stationID, Line: lines[i]})
			}
		}
	}

	return g
}

// FindRouteBFS finds a valid path between startID and endID using BFS
func (g *Graph) FindAllRoutesBFS(startID, endID string, maxRoutes int) ([][]string, [][]string, error) {
	type node struct {
		ID    string
		Path  []string
		Lines []string
	}

	queue := []node{{ID: startID, Path: []string{startID}, Lines: []string{}}}
	var allStations [][]string
	var allLines [][]string
	visited := make(map[string]int) // keeps track of visits to avoid loops

	for len(queue) > 0 && len(allStations) < maxRoutes {
		current := queue[0]
		queue = queue[1:]

		if current.ID == endID {
			stations := make([]string, len(current.Path))
			for i, id := range current.Path {
				stations[i] = g.Names[id]
			}
			allStations = append(allStations, stations)
			allLines = append(allLines, current.Lines)
			continue // do not stop BFS here; keep searching
		}

		// Limit repeated visits to avoid infinite cycles
		key := current.ID + fmt.Sprint(current.Lines)
		if visited[key] >= 1 {
			continue
		}
		visited[key]++

		for _, edge := range g.Nodes[current.ID] {
			nextPath := append([]string{}, current.Path...)
			nextPath = append(nextPath, edge.StationID)

			nextLines := append([]string{}, current.Lines...)
			nextLines = append(nextLines, edge.Line)

			queue = append(queue, node{ID: edge.StationID, Path: nextPath, Lines: nextLines})
		}
	}

	if len(allStations) == 0 {
		return nil, nil, fmt.Errorf("no route found between %s and %s", g.Names[startID], g.Names[endID])
	}

	return allStations, allLines, nil
}
