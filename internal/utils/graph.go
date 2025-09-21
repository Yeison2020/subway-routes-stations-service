package utils


import (
	"fmt"
	"github.com/yeison2020/subway-routing-service/internal/mbta"
	
)


type Edge struct {
	StationID string 
	Line string
}


type Graph struct {
	Nodes map[string][]Edge
	Names map[string]string
}


// BuildGraph builds a bidirectional graph from routes
/*
Each route between stations is shown with two arrows, one in each direction.
*/

func BuildGraph(routes []mbta.Route) *Graph {

	g := &Graph{
		Nodes : make(map[string][]Edge),
		Names : make(map[string]string),
	}


	for _, route := range routes {
		for i, s := range route.Stations {
			g.Names[s.ID] = s.Name

			if i < len(route.Stations) - 1 {
				next := route.Stations[i +1 ]
				g.Nodes[s.ID] = append(g.Nodes[s.ID], Edge{StationID: next.ID, Line: route.Name})
				g.Nodes[next.ID] = append(g.Nodes[next.ID], Edge{StationID: s.ID, Line:route.Name})
			}
		}
	}

	return g

}

// FindRouteBFS finds a valid path between startID and endID using BFS
func (g *Graph) FindRouteBFS(startID, endID string) (stations []string, lines []string, err error) {
	type node struct {
		ID    string
		Path  []string
		Lines []string
	}

	queue := []node{{ID: startID, Path: []string{startID}, Lines: []string{}}}
	visited := make(map[string]bool)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.ID == endID {
			for _, id := range current.Path {
				stations = append(stations, g.Names[id])
			}
			lines = current.Lines
			return stations, lines, nil
		}

		if visited[current.ID] {
			continue
		}
		visited[current.ID] = true

		for _, edge := range g.Nodes[current.ID] {
			if !visited[edge.StationID] {
				nextPath := append([]string{}, current.Path...)
				nextPath = append(nextPath, edge.StationID)

				nextLines := append([]string{}, current.Lines...)
				nextLines = append(nextLines, edge.Line)

				queue = append(queue, node{ID: edge.StationID, Path: nextPath, Lines: nextLines})
			}
		}
	}

	return nil, nil, fmt.Errorf("no route found between %s and %s", startID, endID)
}
