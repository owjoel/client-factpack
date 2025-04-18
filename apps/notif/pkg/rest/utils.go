package rest

import "strconv"

func parseIntWithDefault(str string, defaultVal int) int {
	if value, err := strconv.Atoi(str); err == nil && value > 0 {
		return value
	}
	return defaultVal
}