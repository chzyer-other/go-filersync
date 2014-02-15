package filersync

import (
	"io"
	"os"
	"io/ioutil"
	"testing"
)

func TestStatLinked(t *testing.T) {
	tmpdir, err := ioutil.TempDir("/tmp", "filersync")
	if err != nil { t.Fatal(err) }
	defer os.RemoveAll(tmpdir)

	content := []byte("")
	err = ioutil.WriteFile(tmpdir+"/a.log-1", content, 0666)
	if err != nil { t.Fatal(err) }

	err = ioutil.WriteFile(tmpdir+"/a.log-2", content, 0666)
	if err != nil { t.Fatal(err) }

	sl, err := NewStatLinked([]string {tmpdir + "/a.log-1", tmpdir+"/a.log-2"})
	if err != nil { t.Fatal(err) }

	now, err := sl.First()
	if err != nil {
		t.Fatal("must got one stat")
	}

	next, err := sl.Next(now)
	if err != nil {
		t.Fatal("must got one stat in next")
	}
	if next.Path != tmpdir + "/a.log-2" {
		t.Fatal("next stat not excepted")
	}

	_, err = sl.Prev(now)
	if err != io.EOF {
		t.Fatal("prev of 0 must be eof")
	}

	_, err = sl.Next(next)
	if err != io.EOF {
		t.Fatal("next of last one must be eof")
	}

	prev, err := sl.Prev(next)
	if err != nil {
		t.Fatal("prev of 1 must be not eof")
	}
	if prev.Path != tmpdir + "/a.log-1" {
		t.Fatal("prev of 1 must be alog.-1")
	}
}
