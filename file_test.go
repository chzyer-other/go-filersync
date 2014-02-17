package filersync

import (
	"os"
	"io/ioutil"
	"testing"
)

func TestFile(t *testing.T) {
	tmp, err := ioutil.TempFile("/tmp", "filersync")
	if err != nil { t.Fatal(err) }
	defer os.RemoveAll(tmp.Name())

	stat, err := getStat(tmp.Name())
	if err != nil { t.Fatal(err) }
	_, err = NewFile(nil)
	if err == nil { t.Fatal("must be err") }

	f, err := NewFile(stat)
	if err != nil { t.Fatal(err) }
	
	os.Remove(tmp.Name())
	_, err = NewFile(stat)
	if err == nil { t.Fatal("must be not found file") }

	ino := f.Inode()
	if ino != stat.Inode() { t.Fatal("inode not exist") }
}
