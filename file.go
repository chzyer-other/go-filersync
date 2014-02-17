package filersync

import (
	"os"
	"errors"
)

type File struct {
	*os.File
	Stat *Stat
}

func NewFile(stat *Stat) (f *File, err error) {
	if stat == nil { return nil, errors.New("nil stat") }
	of, err := os.Open(stat.Path)
	if err != nil { return }
	f = &File {
		of,
		stat,
	}
	return
}

func (f *File) Inode() (ino uint64) {
	return f.Stat.Inode()
}

func (f *File) Same(nf *File) bool {
	return f.Stat.SameIno(nf.Stat)
}
