package dispatch_test
import ("testing"; "github.com/wyw14/cry-153/internal/dispatch")
func TestNoDispatchRouteDoesNotBecomeCallableNil(t *testing.T){p:=dispatch.NewPlanner();c:=dispatch.NewCoordinator(p,dispatch.NewState());defer func(){if r:=recover();r!=nil{t.Fatalf("missing route panicked: %v",r)}}();if _,err:=c.Create("flooded");err==nil{t.Fatal("missing route was accepted")}}
