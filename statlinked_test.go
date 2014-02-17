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

	err = ioutil.WriteFile(tmpdir+"/a.log-3", content, 0666)
	if err != nil { t.Fatal(err) }

	sl, err := NewStatLinked([]string {tmpdir + "/a.log-1", tmpdir+"/a.log-2"})
	if err != nil { t.Fatal(err) }

	s1, err := sl.First()
	if err != nil {
		t.Fatal("must got one stat")
	}

	s2, err := sl.Next(s1)
	if err != nil {
		t.Fatal("must got one stat in next")
	}
	if s2.Path != tmpdir + "/a.log-2" {
		t.Fatal("next stat not excepted")
	}

	_, err = sl.Prev(s1)
	if err != io.EOF {
		t.Fatal("prev of 0 must be eof")
	}

	_, err = sl.Next(s2)
	if err != io.EOF {
		t.Fatal("next of last one must be eof")
	}

	ts1, err := sl.Prev(s2)
	if err != nil {
		t.Fatal("prev of 1 must be not eof")
	}
	if ! ts1.SameIno(s1) {
		t.Fatal("prev of 1 must be alog-1")
	}

	// todo delete current stat and try stat
	os.Remove(tmpdir + "/a.log-1")
	os.Remove(tmpdir + "/a.log-2")
	err = sl.UpdatePath([]string {tmpdir+"/a.log-3"})
	if err != nil { t.Fatal(err) }

	s3, err := sl.Next(s2)
	if err != nil { t.Fatal(err) }
	if s3.Path != tmpdir + "/a.log-3" {
		t.Fatal("s2.Next must be s3")
	}

	sl.Remove(s2)
	tts1, err := sl.Prev(s3)
	if err != nil {
		t.Fatal("s3.prev must not be eof")
	}
	if ! tts1.SameIno(s1) {
		t.Fatal("s3.prev must not be s1(s2 deleted)")
	}
}
