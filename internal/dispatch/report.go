package dispatch

import (
	"fmt"
	"time"
)

type Report struct {
	TaskID   string
	Started  time.Time
	Finished time.Time
	Status   string
	Message  string
}

func NewReport(id string) Report { return Report{TaskID: id, Started: time.Now(), Status: "running"} }
func (r *Report) Complete(message string) {
	r.Finished = time.Now()
	r.Status = "complete"
	r.Message = message
}
func (r *Report) Fail(err error) {
	r.Finished = time.Now()
	r.Status = "failed"
	r.Message = err.Error()
}
func (r Report) Duration() time.Duration {
	if r.Finished.IsZero() {
		return time.Since(r.Started)
	}
	return r.Finished.Sub(r.Started)
}
func (r Report) String() string { return fmt.Sprintf("%s %s %s", r.TaskID, r.Status, r.Message) }
