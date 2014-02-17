package filersync

import (
	"errors"
	"regexp"
	"strings"
	"path/filepath"
)

var (
	errMustNotContainWildcard = errors.New("path must not contain wildcard")
	errInvalidPath = errors.New("invalid path")
	LogSuffix []string = []string {"", ".log", `.log-\d+`}
	regexpNumber = regexp.MustCompile(`\d+`)
)

type StringCount struct {
	Content string
	Count int
}

type SameCounter []StringCount
func (s SameCounter) Set(idx int, str string) (exist bool) {
	for idx, i := range s {
		if i.Content == "" { break }
		if i.Content == str {
			i.Count ++
			s[idx] = i
			return true
		}
	}
	s[idx] = StringCount{str, 1}
	return false
}
func (s SameCounter) Max() (string) {
	max := 0
	idx := -1
	for i, a := range s {
		if a.Count > max {
			max = a.Count
			idx = i
		}
	}
	if idx < 0 { return "" }
	return s[idx].Content
}

func extractFileList(path string) (list []string, err error) {
	if strings.Contains(path, "*") {
		err = errMustNotContainWildcard
		return
	}
	if len(path) == 0 {
		err = errInvalidPath
		return
	}
	if path[len(path)-1] != '/' { path += "/" }
	path += "*"
	list, err = filepath.Glob(path)
	return
}

func IsIncludeFileSuffix(path string) bool {
	for _, i := range LogSuffix {
		if strings.HasSuffix(path, i) { return true }
	}
	return false
}

func SelectPath(path string) (filelist []string, err error) {
	fl, err := extractFileList(path)
	if err != nil { return }
	fileList := make(SameCounter, len(fl))
	length := 0
	for _, filepath := range fl {
		filepath = regexpNumber.ReplaceAllString(filepath, `\d+`)
		if ! IsIncludeFileSuffix(filepath) { continue }
		exist := fileList.Set(length, filepath)
		if ! exist { length ++ }
	}
	fileList = fileList[:length]
	if length == 0 { return }

	regexpPath, err := regexp.Compile(fileList.Max())
	if err != nil { return }
	// find match path
	length = 0
	for _, fp := range fl {
		if regexpPath.MatchString(fp) {
			fl[length] = fp
			length += 1
		}
	}
	filelist = fl[:length]
	return
}
