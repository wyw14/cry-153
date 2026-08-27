package plume_test
import ("testing"; "github.com/wyw14/cry-153/internal/plume")
func TestPlumeRetryClassificationDoesNotDependOnLocalizedText(t *testing.T){m:=plume.NewModel("zh");if !plume.RetryableForLocale(m,plume.ErrTemporaryUpstream){t.Fatal("temporary model failure became permanent under localization")}}
