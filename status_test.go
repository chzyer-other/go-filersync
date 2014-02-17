package filersync

import (
	"os"
	"io/ioutil"
	"testing"
)

func TestStatus(t *testing.T) {
	tmpdir, err := ioutil.TempDir("/tmp", "filersync")
	if err != nil { t.Fatal(err) }
	defer os.RemoveAll(tmpdir)

	_, err = NewStatus([]string {"hello"})
	if err == nil { t.Fatal("must be err, not found") }

	content := []byte("hello")
	p1 := tmpdir + "/a.log-1"
	err = ioutil.WriteFile(p1, content, 0666)
	if err != nil { t.Fatal(err) }

	status, err := NewStatus([]string{})
	if err != nil { t.Fatal(err) }

	_, _, err = status.Last()
	if err == nil { t.Fatal("must be error") }

	status.UpdatePath([]string{p1})

	s1, _, err := status.Last()
	if err != nil { t.Fatal(err) }

	_, _, err = status.Prev(s1)
	if err == nil { t.Fatal("must be error") }

	p2 := tmpdir + "/a.log-2"
	err = ioutil.WriteFile(p2, content, 0666)
	if err != nil { t.Fatal(err) }

	status.UpdatePath([]string{p2})
	s2, _, err := status.Next(s1)
	if err != nil { t.Fatal("must not be error", err) }

	ts1, _, err := status.Prev(s2)
	if err != nil { t.Fatal("must not be error", err) }
	if ! ts1.Same(s1) { t.Fatal("s2.prev must be same as s1") }

	_, _, err = status.Next(s2)
	if err == nil { t.Fatal("must be error") }

	finishs1 := false
	change := func() {
		finishs1 = true
	}
	status.UpdateOffset(s1, 2)

	offset, err := status.Offset(s1)
	if err != nil { t.Fatal(err) }
	if offset != 2 { t.Fatal("s1.offset must be 2") }
	if status.IsFinish(s1) { t.Fatal("this must not finish") }

	status.UpdateOffset(s1, 5)
	if ! status.IsFinish(s1) { t.Fatal("this must finish") }
	status.SetNotifyFunc(change)

	status.UpdateOffset(s1, 5)
	if ! finishs1 { t.Fatal("this must be execute") }
	finishs1 = false

	p3 := tmpdir + "/a.log-3"
	err = ioutil.WriteFile(p3, content, 0666)
	if err != nil { t.Fatal(err) }
	status.UpdatePath([]string{p3})
	if ! finishs1 { t.Fatal("onchange must be trigger") }
}
