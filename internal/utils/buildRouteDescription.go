package utils

import "fmt"

func BuildRouteDescriptions(allStations [][]string, allLines [][]string) []string {
	descriptions := make([]string, len(allStations))

	for i := range allStations {
		stations := allStations[i]
		lines := allLines[i]

		if len(stations) == 0 || len(lines) == 0 {
			descriptions[i] = ""
			continue
		}

		desc := "Start at " + stations[0]
		currentLine := lines[0]

		for j := 1; j < len(stations); j++ {
			if lines[j-1] != currentLine {
				desc += fmt.Sprintf(", transfer at %s to %s", stations[j-1], lines[j-1])
				currentLine = lines[j-1]
			}
		}

		desc += fmt.Sprintf(", take %s to %s.", currentLine, stations[len(stations)-1])
		descriptions[i] = desc
	}

	return descriptions
}