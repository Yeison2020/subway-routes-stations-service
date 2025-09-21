package utils

import "fmt"

func BuildRouteDescription(stations, lines []string) string {
	if len(stations) == 0 || len(lines) == 0 {
		return ""
	}

	description := "Start at " + stations[0]
	currentLine := lines[0]

	for i := 1; i < len(stations); i++ {
		if lines[i-1] != currentLine {
			description += fmt.Sprintf(", transfer at %s to %s", stations[i-1], lines[i-1])
			currentLine = lines[i-1]
		}
	}

	description += fmt.Sprintf(", take %s to %s.", currentLine, stations[len(stations)-1])
	return description
}
