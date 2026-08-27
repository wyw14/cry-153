package gate_test
import ("testing"; "time"; "github.com/wyw14/cry-153/internal/gate")
func TestGateBackoffSaturatesWithoutDurationOverflow(t *testing.T){max:=8*time.Second;if got:=gate.BoundedBackoff(100,max);got!=max{t.Fatalf("large retry overflowed to %s",got)};for i:=0;i<16;i++{got:=gate.BoundedBackoff(i,max);if got<0||got>max{t.Fatalf("invalid delay %d: %s",i,got)}}}
