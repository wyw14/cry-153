package incident_test
import ("testing"; "github.com/wyw14/cry-153/internal/alert"; "github.com/wyw14/cry-153/internal/incident"; "github.com/wyw14/cry-153/internal/journal")
func TestIncidentCloseAllowsBoundedFinalAlertPublication(t *testing.T){st:=incident.NewState();a:=alert.NewPublisher(alert.NewState());c:=incident.NewCoordinator(st,journal.NewStore(""),a);item:=c.Open("spillway");if err:=c.Close(item.ID);err!=nil{t.Fatal(err)};alerts:=a.Sent();if len(alerts)!=1||alerts[0].Active{t.Fatalf("final inactive alert missing: %#v",alerts)}}
