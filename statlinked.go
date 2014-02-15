package filersync

import (
	"io"
	"sync"
	// "os"
)

type StatLinked struct {
	stat []*Stat
	statL sync.RWMutex
}

func NewStatLinked(paths []string) (sl *StatLinked, err error) {
	sl = &StatLinked {
	}
	err = sl.UpdatePaths(paths)
	if err != nil { return }
	return
}

func (sl *StatLinked) UpdatePaths(paths []string) (err error) {
	stats := make(StatSlice, len(paths))
	for idx, path := range paths {
		stats[idx], err = getStat(path)
		if err != nil { return }
	}
	stats.Sort()

	sl.statL.Lock()
	defer sl.statL.Unlock()
	sl.stat = stats
	return
}

func (sl *StatLinked) Prev(s *Stat) (ns *Stat, err error) {
	sl.statL.RLock()
	defer sl.statL.RUnlock()

	err = io.EOF
	for idx, i := range sl.stat {
		if i.Inode() == s.Inode() || s.MBefore(i) {
			if idx > 0 { return sl.stat[idx-1], nil }
			return
		}
	}
	return
}

func (sl *StatLinked) Next(s *Stat) (ns *Stat, err error) {
	sl.statL.RLock()
	defer sl.statL.RUnlock()

	err = io.EOF
	for idx, i := range sl.stat {
		if i.Inode() == s.Inode() {
			if len(sl.stat) > idx + 1 {
				return sl.stat[idx+1], nil
			}
			return
		}
		if s.MBefore(i) { return i, nil }
	}
	err = io.EOF
	return
}

func (sl *StatLinked) First() (s *Stat, err error) {
	sl.statL.RLock()
	defer sl.statL.RUnlock()

	err = io.EOF
	if len(sl.stat) == 0 { return }
	return sl.stat[0], nil
}

func (sl *StatLinked) Last() (s *Stat, err error) {
	sl.statL.RLock()
	defer sl.statL.RUnlock()

	err = io.EOF
	if len(sl.stat) == 0 { return }
	return sl.stat[len(sl.stat)-1], nil
}
