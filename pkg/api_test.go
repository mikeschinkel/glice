package glice

//func TestGitHubAPINoKey(t *testing.T) {
//	c := context.Background()
//	d := &Dependency{
//		Import:    "github.com/ribice/kiss",
//		Host:    "github.com",
//		Author:  "ribice",
//		Project: "kiss",
//	}
//
//	gc := NewHostClient(c, map[string]string{}, false)
//	err := gc.GetDependencyLicense(c, d)
//	if err != nil {
//		t.Error(err)
//	}
//
//	if d.Shortname != color.New(color.FgGreen).Sprintf("MIT") {
//		t.Errorf("API did not return correct license or color.")
//	}
//
//}
//
//func TestNonexistentLicense(t *testing.T) {
//
//	c := context.Background()
//	d := &Dependency{
//		Import:    "github.com/denysdovhan/wtfjs",
//		Host:    "github.com",
//		Author:  "denysdovhan",
//		Project: "wtfjs",
//	}
//
//	gc := NewHostClient(c, map[string]string{}, false)
//	err := gc.GetDependencyLicense(c, d)
//	if err != nil {
//		t.Error(err)
//	}
//
//	if d.Shortname != color.New(color.FgYellow).Sprintf("wtfpl") {
//		t.Errorf("API did not return correct license or color.")
//	}
//
//}
//
//func TestGitHubAPIWithKey(t *testing.T) {
//
//	c := context.Background()
//	d := &Dependency{
//		Import:    "github.com/ribice/kiss",
//		Host:    "github.com",
//		Author:  "ribice",
//		Project: "kiss",
//	}
//
//	v := map[string]string{
//		"github.com": "apikey",
//	}
//
//	gc := NewHostClient(c, v, false)
//	err := gc.GetDependencyLicense(c, d)
//	if err == nil {
//		t.Error("expected bad credentials error")
//	}
//
//}
//
//func TestGitHubAPIWithKeyAndThanks(t *testing.T) {
//
//	c := context.Background()
//	d := &Dependency{
//		Import:    "github.com/ribice/kiss",
//		Host:    "github.com",
//		Author:  "ribice",
//		Project: "kiss",
//	}
//
//	v := map[string]string{
//		"github.com": "apikey",
//	}
//
//	gc := NewHostClient(c, v, true)
//
//	err := gc.GetDependencyLicense(c, d)
//	if err == nil {
//		t.Error("expected bad credentials error")
//	}
//
//}
