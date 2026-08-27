package service

import (
	"fmt"
	"runtime"
	"time"
)

type Diagnostic struct {
	At         time.Time
	GoVersion  string
	Goroutines int
	Status     string
}

func CollectDiagnostic() Diagnostic {
	return Diagnostic{At: time.Now(), GoVersion: runtime.Version(), Goroutines: runtime.NumGoroutine(), Status: "ok"}
}
func (d Diagnostic) String() string {
	return fmt.Sprintf("%s %s goroutines=%d", d.Status, d.GoVersion, d.Goroutines)
}
func (d Diagnostic) Healthy() bool { return d.Status == "ok" && d.Goroutines > 0 }
func MergeDiagnostics(a, b Diagnostic) Diagnostic {
	if b.At.After(a.At) {
		return b
	}
	return a
}
