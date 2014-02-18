package filersync

import (
	"os"
	"testing"
	"io/ioutil"
)

func TestFileSelect(t *testing.T) {
	tmpdir, err := ioutil.TempDir("/tmp", "filersync")
	if err != nil { t.Fatal(err) }
	defer os.RemoveAll(tmpdir)

	p1, p2 := tmpdir + "/a.log-1", tmpdir + "/a.log-2"
	err = ioutil.WriteFile(p1, []byte("hello, this is p1"), 0666)
	if err != nil { t.Fatal(err) }
	err = ioutil.WriteFile(p2, []byte("hello, this is p2"), 0666)
	if err != nil { t.Fatal(err) }

	status, err := NewStatus([]string {})
	if err != nil { t.Fatal(err) }
	
	fselect, err := NewFileSelector(status)
	if err != nil { t.Fatal(err) }

	_, _, err = fselect.GetNewFile()
	if err == nil { t.Fatal("must be error") }

	status.UpdatePath([]string{p1})

	f, offset, err := fselect.GetNewFile()
	if err != nil { t.Fatal(err) }
	if offset != 0 { t.Fatal("offset must be 0") }

	status.UpdateOffset(f, 2)
	f, offset, err = fselect.GetNewFile()
	if err != nil { t.Fatal(err) }
	if offset != 2 { t.Fatal("offset must be 2") }

	_, _, err = fselect.GetOldFile()
	if err == nil { t.Fatal("must be error") }

	status.UpdatePath([]string{p2})
	f1, _, err := fselect.GetOldFile()
	if err != nil { t.Fatal("must not be error") }
	if f1.Name() != p1 { t.Fatal("f1.path must be", p1, "got", f1.Name()) }

	f2, _, err := fselect.GetNewFile()
	if err != nil { t.Fatal("must not be error") }
	if f2.Name() != p2 { t.Fatal("f2.path must be", p2, "got", f2.Name()) }
}
