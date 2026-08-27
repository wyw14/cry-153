package incident

import (
	"fmt"
	"github.com/wyw14/cry-153/internal/model"
)

func Transition(item model.Incident, next model.IncidentState) (model.Incident, error) {
	if item.State == model.StateClosed && next != model.StateClosed {
		return item, fmt.Errorf("closed incident cannot transition")
	}
	item.State = next
	return item, nil
}
func Label(item model.Incident) string { return fmt.Sprintf("%s:%s", item.ID, item.State) }
