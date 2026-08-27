package dispatch_test
import ("testing"; "time"; "github.com/wyw14/cry-153/internal/dispatch")
func TestUpdatedPlumeDeadlineReordersDispatchHeap(t *testing.T){q:=dispatch.NewScheduler();q.Schedule(dispatch.Task{ID:"routine",Priority:1,Deadline:time.Now().Add(time.Hour)});q.Schedule(dispatch.Task{ID:"urgent",Priority:2,Deadline:time.Now().Add(2*time.Hour)});q.Update("routine",10,time.Now().Add(time.Minute));next,ok:=q.Next();if !ok||next.ID!="routine"{t.Fatalf("updated task was not promoted: %#v",next)}}
