package filersync

import (
	"os"
	"time"
	"testing"
	"io/ioutil"
)

func TestFileSync(t *testing.T) {
	tmpdir, err := ioutil.TempDir("/tmp", "filersync-filesync")
	if err != nil { t.Fatal(err) }
	defer os.RemoveAll(tmpdir)

	fs, _ := NewFileSync(tmpdir)
	time.Sleep(100*time.Millisecond)

	content := []byte("hello\n1\n")
	p1, p2 := tmpdir+"/a.log-1", tmpdir+"/a.log-2"
	p3, p4 := tmpdir+"/a.log-3", tmpdir+"/a.log-4"
	err = ioutil.WriteFile(p1, content, 0666)
	if err != nil { t.Fatal(err) }

	time.Sleep(100*time.Millisecond)
	c, path, err := fs.Readline()
	if err != nil { t.Fatal("must not be error") }
	if c != "hello" { t.Fatal("content msut be `hello` got ", c) }
	if path != p1 { t.Fatal("path must be p1") }

	err = ioutil.WriteFile(p2, []byte("h2\n2\n"), 0666)
	if err != nil { t.Fatal(err) }

	time.Sleep(200*time.Millisecond)

	c, path, err = fs.Readline()
	if err != nil { t.Fatal("content must not be null") }
	if c != "h2" { t.Fatal("content must be `h2`", c) }
	if path != p2 { t.Fatal("path must be p2", path) }

	c, path, err = fs.Readline()
	if err != nil { t.Fatal("readline must not be error") }
	if c != "2" { t.Fatal("content must be `2`") }
	if path != p2 { t.Fatal("path must be p2", path) }

	c, path, err = fs.Readline()
	if err != nil { t.Fatal("readline must not be error") }
	if c != "1" { t.Fatal("content must be `1`") }
	if path != p1 { t.Fatal("path must be p1", path) }

	_, _, err = fs.Readline()
	if err == nil { t.Fatal("this must be eof") }

	err = ioutil.WriteFile(p3, []byte("h3\nh3-2"), 0666)
	if err != nil { t.Fatal(err) }

	err = ioutil.WriteFile(p4, []byte("h4\nh4-2"), 0666)
	if err != nil { t.Fatal(err) }

	time.Sleep(200*time.Millisecond)

	c, path, err = fs.Readline()
	if err != nil { t.Fatal(err) }
	if path != p4 { t.Fatal("path must be p4", path) }
	if c != "h3-2h4" { t.Fatal("content must be h3-2h4, got", c) }

	c, path, err = fs.Readline()
	if err != nil { t.Fatal(err) }
	if path != p3 { t.Fatal("path msut be p3") }
	if c != "h3" { t.Fatal("content msut be h3") }

	_, _, err = fs.Readline()
}
