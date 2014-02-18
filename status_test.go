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

	_, err = status.Last()
	if err == nil { t.Fatal("must be error") }

	status.UpdatePath([]string{p1})

	s1, err := status.Last()
	if err != nil { t.Fatal(err) }

	_, err = status.Prev(s1)
	if err == nil { t.Fatal("must be error") }

	p2 := tmpdir + "/a.log-2"
	err = ioutil.WriteFile(p2, content, 0666)
	if err != nil { t.Fatal(err) }

	status.UpdatePath([]string{p2})
	s2, err := status.Next(s1)
	if err != nil { t.Fatal("must not be error", err) }

	ts1, err := status.Prev(s2)
	if err != nil { t.Fatal("must not be error", err) }
	if ! ts1.Same(s1) { t.Fatal("s2.prev must be same as s1") }

	_, err = status.Next(s2)
	if err == nil { t.Fatal("must be error") }

	finishs1 := false
	change := func() {
		finishs1 = true
	}
	status.UpdateOffset(s1, 2)

	offset := status.Offset(s1)
	if offset != 2 { t.Fatal("s1.offset must be 2, offset", offset) }
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

	offset = status.Offset(s2)
	if offset != 0 { t.Fatal("s2.offset must be 0", offset) }
	f3, err := status.Next(s2)
	if err != nil { t.Fatal("s2.next must be s3") }
	of3 := status.Offset(f3)
	if of3 != 0 { t.Fatal("f3.offset must be 0", offset) }

	err = ioutil.WriteFile(p3, append(content, '1'), 0666)
	if err != nil { t.Fatal(err) }
	status.UpdatePath([]string{p3}) // trigger increase
	of3 = status.Offset(f3)
	status.Sync(f3)
	if f3.Stat.Size != int64(len(content) + 1) { t.Fatal("content length must be", len(content)+1, "got", f3.Stat.Size) }

	status.UpdateOffset(f3, 5)
	if status.IsFinish(f3) { t.Fatal("must not finish") }
	status.UpdateOffset(f3, 6)
	if ! status.IsFinish(f3) { t.Fatal("must finish") }

	err = ioutil.WriteFile(p3, append(content, '1', '2'), 0666)
	if err != nil { t.Fatal(err) }
	status.UpdatePath([]string{p3}) // trigger increase

	if status.IsFinish(f3) { t.Fatal("must not finish") }
}
