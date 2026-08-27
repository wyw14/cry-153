package station

import (
	"fmt"
	"github.com/wyw14/cry-153/internal/model"
)

func ValidateSource(source Source) error {
	if source.ID == "" {
		return fmt.Errorf("source id required")
	}
	if len(source.Values) == 0 {
		return fmt.Errorf("source has no values")
	}
	for _, v := range source.Values {
		if v < 0 {
			return fmt.Errorf("negative sample")
		}
	}
	return nil
}
func ToBatch(id string, readings []model.Reading) model.SampleBatch {
	return model.SampleBatch{IncidentID: id, Readings: model.CloneReadings(readings)}
}
func OnlineCount(status map[string]bool) int {
	n := 0
	for _, v := range status {
		if v {
			n++
		}
	}
	return n
}
