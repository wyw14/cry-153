package sample

import (
	"fmt"
	"github.com/wyw14/cry-153/internal/model"
)

func ValidateBatch(batch model.SampleBatch) error {
	if batch.IncidentID == "" {
		return fmt.Errorf("incident required")
	}
	if len(batch.Readings) == 0 {
		return fmt.Errorf("readings required")
	}
	for _, r := range batch.Readings {
		if err := model.ValidateReading(r); err != nil {
			return err
		}
	}
	return nil
}
func Normalize(batch model.SampleBatch) model.SampleBatch {
	batch.Readings = model.CloneReadings(batch.Readings)
	return batch
}
func CountByStation(readings []model.Reading) map[string]int {
	out := map[string]int{}
	for _, r := range readings {
		out[r.StationID]++
	}
	return out
}
func MaxConcentration(readings []model.Reading) float64 {
	max := 0.0
	for _, r := range readings {
		if r.Concentration > max {
			max = r.Concentration
		}
	}
	return max
}
