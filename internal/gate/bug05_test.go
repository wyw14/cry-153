package gate_test
import ("context"; "testing"; "github.com/wyw14/cry-153/internal/gate")
func TestCanceledGateAcquireCannotInflateInFlightCapacity(t *testing.T){cap:=gate.NewCapacity(1);c:=gate.NewClient("http://127.0.0.1:1",cap);ctx,cancel:=context.WithCancel(context.Background());cancel();for i:=0;i<8;i++{if err:=c.Close(ctx,"north","inc");err==nil{t.Fatal("canceled close unexpectedly succeeded")}};if got:=cap.InFlight();got!=0||got>cap.Limit(){t.Fatalf("capacity accounting inflated: %d",got)}}
