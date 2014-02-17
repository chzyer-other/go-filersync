package filersync

import (
)

type FileSelector struct {
	status *Status
}

func NewFileSelector(status *Status) (fs *FileSelect, err error) {
	fs = &FileSelector {
		status: status,
	}
	return
}

func (fs *FileSelector) OnFileChange() {
}

func (fs *FileSelector) GetNewFile() (f *File, offset int64) {
	return
}

func (fs *FileSelector) GetOldFile() (f *File, offset int64) {
	return
}
