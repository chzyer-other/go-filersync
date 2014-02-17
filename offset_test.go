package filersync

import (
	"os"
	"io/ioutil"
	"testing"
)

func TestOffset(t *testing.T) {
	tmpdir, err := ioutil.TempDir("/tmp", "filersync")
	if err != nil { t.Fatal(err) }
	defer os.RemoveAll(tmpdir)
	content := []byte("abcefg")

	err = ioutil.WriteFile(tmpdir+"/a.log-1", content, 0666)
	if err != nil { t.Fatal(err) }
	defer os.Remove(tmpdir + "/a.log-1")

	err = ioutil.WriteFile(tmpdir+"/a.log-2", content, 0666)
	if err != nil { t.Fatal(err) }
	defer os.Remove(tmpdir + "/a.log-2")
	
	stat1, err := getStat(tmpdir+"/a.log-1")
	if err != nil { t.Fatal(err) }

	stat2, err := getStat(tmpdir + "/a.log-2")
	if err != nil { t.Fatal(err) }
	
	o := NewOffset()
	_, ok := o.GetStatOffset(stat1)
	if ok { t.Fatal("must not be ok") }

	o.SetStatOffset(stat1, 1)

	s, ok := o.GetStatOffset(stat1)
	if ! ok { t.Fatal("must be ok") }
	if s != 1 { t.Fatal("s1.offset must be 1") }

	finish := o.IsFinish(stat1)
	if finish { t.Fatal("must not be finish") }

	o.AddStatOffset(stat1, 2)
	s, ok = o.GetStatOffset(stat1)
	if ! ok { t.Fatal("must be ok") }
	if s != 3 { t.Fatal("s1.offset must be 3") }

	finish = o.IsFinish(stat1)
	if finish { t.Fatal("must not be finish") }

	o.SetStatLimit(stat1, 3)
	finish = o.IsFinish(stat1)
	if ! finish { t.Fatal("must be finish") }

	ok = o.Remove(stat1)
	if ! ok { t.Fatal("must can remove") }
	ok = o.Remove(stat1)
	if ok { t.Fatal("must could not remove") }

	_, ok = o.GetStatOffset(stat1)
	if ok { t.Fatal("must not be ok") }

	ok = o.IsFinish(stat1)
	if ok { t.Fatal("must not finish") }

	o.AddStatOffset(stat2, 3)
	s, ok = o.GetStatOffset(stat2)
	if ! ok { t.Fatal("must be ok") }
	if s != 3 { t.Fatal("must be 3") }

}
