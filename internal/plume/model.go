package plume

import (
	"errors"
	"fmt"
	"github.com/wyw14/cry-153/internal/model"
)

var ErrTemporaryUpstream = errors.New("temporary upstream gap")
var ErrPermanentModel = errors.New("permanent model failure")

type Localizer interface{ Localize(error) string }
type Model struct{ locale string }

func NewModel(locale string) *Model { return &Model{locale: locale} }
func (m *Model) Estimate(readings []model.Reading) (float64, error) {
	if len(readings) == 0 {
		return 0, ErrTemporaryUpstream
	}
	for _, r := range readings {
		if r.Concentration < 0 {
			return 0, ErrPermanentModel
		}
	}
	return readings[len(readings)-1].Concentration * 1.25, nil
}
func (m *Model) Display(err error) string {
	if err == nil {
		return ""
	}
	if m.locale == "zh" {
		return "上游暂时缺口"
	}
	return fmt.Sprint(err)
}
