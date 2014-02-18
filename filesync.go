package filersync

import (
	"io"
	"bytes"
	"errors"
)

var (
	exceptLen = 10240
)

type FileSync struct {
	newPathChan chan []string
	status *Status
	fselect *FileSelector
}

func NewFileSync(path string) (fs *FileSync, err error) {
	status, err := NewStatus(nil)
	if err != nil { return }

	fselect, err := NewFileSelector(status)
	if err != nil { return }

	fs = &FileSync {
		newPathChan: KeepReturnNewPath(path),
		status: status,
		fselect: fselect,
	}
	go fs.WaitForUpdatePath()
	return
}

func (fs *FileSync) WaitForUpdatePath() {
	for ps := range fs.newPathChan {
		fs.status.UpdatePath(ps)
	}
}

func (fs *FileSync) Readline() (data, path string, err error) {
	var d []byte
	f, offset, err := fs.fselect.GetNewFile()
	if err == nil {
		d, err = fs.readline(f, offset)
		if err == io.EOF {
		} else if err != nil {
			return
		}
	}
	if err != nil {
		f, offset, err = fs.fselect.GetOldFile()
		if err != nil { return }
		d, err = fs.readline(f, offset)
		if err != nil { return }
	}
	data = string(d)
	fs.status.AddOffset(f, int64(len(d)+1))
	path = f.Name()
	return
}

func (fs *FileSync) getFileLastLine(f *File) (data []byte, err error) {
	tmpret := make([]byte, exceptLen)
	f.Seek(-int64(exceptLen), 2)
	n, err := f.Read(tmpret)
	if err == nil {
		tmpret = tmpret[:n]
		idx := bytes.LastIndex(tmpret, []byte("\n"))
		if idx >= 0 && idx < n-1 {
			return tmpret[idx+1:], nil
		}
	}
	err = errors.New("last line not found")
	return
}

func (fs *FileSync) readline(f *File, offset int64) (data []byte, err error) {
	ret := make([]byte, exceptLen)
	f.Seek(offset, 0)
	n, err := f.Read(ret)
	if err != nil { return }
	ret = ret[:n]

	idx := bytes.IndexByte(ret, '\n')
	if idx < 0 {
		err = io.EOF
		return
	}
	data = ret[:idx]

	prevfile, e := fs.status.Prev(f)
	if e == nil {
		pret, e := fs.getFileLastLine(prevfile)
		if len(data) > 0 && e == nil {
			data = append(pret, data...)
		}
	}

	return
}
