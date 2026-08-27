package alert

import "strings"

func Summary(active bool, intakes []string) string {
	state := "inactive"
	if active {
		state = "active"
	}
	return state + " alerts for " + strings.Join(intakes, ",")
}
