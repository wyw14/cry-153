package timeline_test

import (
	"github.com/wyw14/cry-153/internal/station"
	"github.com/wyw14/cry-153/internal/timeline"
	"sync"
	"testing"
	"time"
)

func TestTimelineSnapshotDoesNotIterateLiveStationMap(t *testing.T) {
	reg := station.NewRegistry()
	w := timeline.NewWindow()
	c := timeline.NewCoordinator(w, reg)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 500; i++ {
			reg.AddReading(station.FreshReading("north", float64(i)))
			c.Record("inc", "sample", "reading")
		}
	}()
	for i := 0; i < 500; i++ {
		snapshot := c.StationSnapshot()
		for id, items := range snapshot {
			_ = id
			_ = len(items)
		}
		_ = reg.Status()
		time.Sleep(time.Microsecond)
	}
	wg.Wait()
}
