package filersync

import (
	"os"
	"time"
	"testing"
	"io/ioutil"
)
var _ = os.Remove
func init() {
	UpgradeInterval = 100 * time.Millisecond
}

func TestUpgrade(t *testing.T) {
	tmpdir, err := ioutil.TempDir("/tmp", "filersync-upgrade")
	if err != nil { t.Fatal(err) }
	defer os.RemoveAll(tmpdir)

	err = ioutil.WriteFile(tmpdir + "/a.log-44", []byte("helo"), 0666)
	if err != nil { t.Fatal(err) }
	ch := KeepReturnNewPath(tmpdir)
	select {
	case <-time.After(100*time.Millisecond):
		t.Fatal("get path timeout")
	case a:= <- ch :
		if a[0] != tmpdir + "/a.log-44" {
			t.Fatal("result not excepted", a)
		}
	}
	time.Sleep(200 * time.Millisecond)

	select {
	case <-time.After(100*time.Millisecond):
	case a := <- ch:
		t.Fatal("may not return data", a)
	}

	err = ioutil.WriteFile(tmpdir + "/a.log-131", []byte("a"), 0666)
	if err != nil { t.Fatal(err) }

	time.Sleep(200 * time.Millisecond)

	select {
	case <-time.After(100*time.Millisecond):
		t.Fatal("timeout")
	case a := <- ch:
		if len(a) != 2 {
			t.Fatal("len of result may be 2, got", len(a))
		}
		if ! inArray(tmpdir + "/a.log-131", a) || ! inArray(tmpdir + "/a.log-44", a) {
			t.Fatal("result not excepted", a)
		}
	}
}
