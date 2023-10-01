package discover

import (
	"os"
	"testing"

	"github.com/chenjie199234/Corelib/discover"
)

func Test_Discover(t *testing.T) {
	os.Setenv("ADMIN_SERVICE_PROJECT", "corp")
	os.Setenv("ADMIN_SERVICE_GROUP", "system")
	os.Setenv("ADMIN_SERVICE_WEB_HOST", "localhost")
	os.Setenv("ADMIN_SERVICE_WEB_PORT", "8000")
	os.Setenv("ADMIN_SERVICE_DISCOVER_ACCESS_KEY", "accesskey_watch_discover")
	di, e := NewAdminDiscover("corp", "system", "test", "corp", "system", "admin", nil)
	if e != nil {
		t.Fatal(e)
		return
	}
	notice, cancel := di.GetNotice()
	for {
		_, ok := <-notice
		if !ok {
			t.Fatal("stopped")
			cancel()
			return
		}
		addrs, version, e := di.GetAddrs(discover.NotNeed)
		if e != nil {
			t.Fatal(e)
			cancel()
			return
		}
		t.Log(addrs)
		t.Log(version)
		t.Log("//////////////////////////////////")
	}
}
