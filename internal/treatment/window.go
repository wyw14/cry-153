package treatment

import (
	"github.com/wyw14/cry-153/internal/model"
	"sort"
	"time"
)

func CopyWindow(window model.SampleWindow) model.SampleWindow {
	window.Readings = model.CloneReadings(window.Readings)
	return window
}
func SortWindow(window model.SampleWindow) model.SampleWindow {
	window = CopyWindow(window)
	sort.Slice(window.Readings, func(i, j int) bool { return window.Readings[i].Concentration < window.Readings[j].Concentration })
	return window
}
func WindowAge(window model.SampleWindow, now time.Time) time.Duration {
	if len(window.Readings) == 0 {
		return 0
	}
	return now.Sub(window.Readings[0].At)
}
func Within(window model.SampleWindow, span time.Duration, now time.Time) bool {
	return WindowAge(window, now) <= span
}
