package intake_test
import ("context"; "sync"; "testing"; "time"; "github.com/wyw14/cry-153/internal/intake")
type driver struct{}
func(driver)Close(context.Context,string,string)error{time.Sleep(time.Millisecond);return nil}
func TestOverlappingIntakeClosuresUseGlobalLockOrder(t *testing.T){locks:=intake.NewLocks([]string{"north","south"});c:=intake.NewCoordinator(locks,driver{});ctx,cancel:=context.WithTimeout(context.Background(),300*time.Millisecond);defer cancel();var wg sync.WaitGroup;errs:=make(chan error,2);for _,ids:=range [][]string{{"north","south"},{"south","north"}}{wg.Add(1);go func(x []string){defer wg.Done();errs<-c.CloseGroup(ctx,"inc",x)}(ids)};done:=make(chan struct{});go func(){wg.Wait();close(done)}();select{case <-done:case <-ctx.Done():t.Fatal("overlapping closures deadlocked")};close(errs);for err:=range errs{if err!=nil{t.Fatal(err)}}}
