package filersync

import (
	"os"
	"testing"
	"io/ioutil"
)


func TestStatSelector(t *testing.T) {
	tmpdir, err := ioutil.TempDir("/tmp", "filersync")
	if err != nil { t.Fatal(err) }
	defer os.RemoveAll(tmpdir)
	content := []byte("")

	err = ioutil.WriteFile(tmpdir+"/a.log-1", content, 0666)
	if err != nil { t.Fatal(err) }

	paths, err := SelectPath(tmpdir)
	if err != nil { t.Fatal(err) }
	if ! inArray(tmpdir+"/a.log-1", paths) {
		t.Fatal("a.log not in array")
	}

	err = ioutil.WriteFile(tmpdir+"/a.log-52347", content, 0666)
	if err != nil { t.Fatal(err) }
	paths, err = SelectPath(tmpdir)
	if err != nil { t.Fatal(err) }
	if ! inArray(tmpdir+"/a.log-1", paths) || ! inArray(tmpdir+"/a.log-52347", paths) {
		t.Fatal("result not except")
	}

	os.Remove(tmpdir+"/a.log-1")
	os.Remove(tmpdir+"/a.log-52347")
	paths, err = SelectPath(tmpdir)
	if err != nil { t.Fatal(err) }
	if len(paths) != 0 {
		t.Fatal("len must be 0")
	}
}
