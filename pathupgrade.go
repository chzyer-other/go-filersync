package filersync

import (
	"time"
)

var (
	UpgradeInterval = time.Second
)

func KeepReturnNewPath(path string) (pathchan chan []string) {
	pathchan = make(chan []string)
	go keepReturnNewpath(path, pathchan)
	return
}

func keepReturnNewpath(path string, pathchan chan []string) {
	for {
		fileList, err := SelectPath(path)
		if err != nil { panic(err) }
		if fileListChanged(fileList) {
			pathchan <- fileList
		}
		time.Sleep(UpgradeInterval)
	}
}

var tmpFileList []string
func fileListChanged(fileList []string) bool {
	if tmpFileList == nil {
		tmpFileList = fileList
		return true
	}

	if len(tmpFileList) != len(fileList) {
		tmpFileList = fileList
		return true
	}

	for _, f := range tmpFileList {
		if ! inArray(f, fileList) {
			tmpFileList = fileList
			return true
		}
	}
	return false
}

func inArray(a string, b []string) bool {
	for _, c := range b { if c == a { return true } }
	return false
}
