package filersync

import (
	"io"
	"sync"
)


type FileSelector struct {
	status *Status
	lastFile *File
	prevFile *File
	lastFileL, prevFileL sync.Mutex
}

func NewFileSelector(status *Status) (fs *FileSelector, err error) {
	fs = &FileSelector {
		status: status,
	}
	status.SetNotifyFunc(fs.OnFileChange)
	return
}

func (fs *FileSelector) OnFileChange() {
	fs.lastFileL.Lock()
	fs.lastFile = nil
	fs.lastFileL.Unlock()

	fs.prevFileL.Lock()
	fs.prevFile = nil
	fs.prevFileL.Unlock()
}

func (fs *FileSelector) GetNewFile() (f *File, offset int64, err error) {
	fs.lastFileL.Lock()
	defer fs.lastFileL.Unlock()
	alreadyLast := false
	if fs.lastFile == nil {
		fs.lastFile, err = fs.status.Last()
		if err != nil { return }
		alreadyLast = true
	}
	if fs.status.IsFinish(fs.lastFile) {
		if alreadyLast { return nil, 0, io.EOF }
		tf, err := fs.status.Last()
		if err != nil { return nil, 0, err }
		if fs.status.IsFinish(tf) { return nil, 0, io.EOF }
		fs.lastFile = tf
	}
	f, offset = fs.lastFile, fs.status.Offset(fs.lastFile)
	fs.status.Sync(fs.lastFile)
	return
}

func (fs *FileSelector) GetOldFile() (f *File, offset int64, err error) {
	fs.prevFileL.Lock()
	defer fs.prevFileL.Unlock()

	if fs.prevFile != nil {
		for fs.status.IsFinish(fs.lastFile) {
			fs.prevFile, err = fs.status.Prev(fs.prevFile)
			if err != nil { return } // at first file
		}
		f, offset = fs.prevFile, fs.status.Offset(fs.prevFile)
		return
	}

	lastFile := fs.lastFile
	if lastFile == nil {
		lastFile, err = fs.status.Last()
		if err != nil { return }
	}
	fs.prevFile, err = fs.status.Prev(lastFile)
	if err != nil { return }
	for fs.status.IsFinish(fs.prevFile) {
		fs.prevFile, err = fs.status.Prev(fs.prevFile)
		if err != nil { return }
	}
	f, offset = fs.prevFile, fs.status.Offset(fs.prevFile)
	return
}
