package filersync

import (
	"io"
	"sync"
	// "os"
)

type StatLinked struct {
	stat StatSlice
	statL sync.RWMutex
	addfunc func()
	increasefunc func(*Stat, int64)
}

func NewStatLinked(paths []string) (sl *StatLinked, err error) {
	sl = &StatLinked {
	}
	err = sl.UpdatePath(paths)
	if err != nil { return }
	return
}

func (sl *StatLinked) SetOnSizeIncrease(f func(*Stat, int64)) {
	sl.increasefunc = f
}

func (sl *StatLinked) search(s *Stat) (idx int) {
	for idx, i := range sl.stat {
		if i.Inode() == s.Inode() { return idx }
	}
	return -1
}

func (sl *StatLinked) Search(s *Stat) (idx int) {
	sl.statL.RLock()
	defer sl.statL.RUnlock()
	return sl.search(s)
}

type Increase struct {
	stat *Stat
	size int64
}

func (sl *StatLinked) UpdatePath(paths []string) (err error) {
	stats := make(StatSlice, len(paths))
	for idx, path := range paths {
		stats[idx], err = getStat(path)
		if err != nil { return }
	}

	needNotify := false
	increase := make([]Increase, len(stats))
	length := 0
	sl.statL.Lock()
	for _, s := range stats {
		idx := sl.search(s)
		if idx >= 0 {
			olds := sl.stat[idx]
			if olds.Size < s.Size {
				increase[length] = Increase{s, s.Size}
			}
			sl.stat[idx] = s
		} else {
			sl.stat = append(sl.stat, s)
			needNotify = true
		}
	}
	sl.stat.Sort()
	sl.statL.Unlock()
	if length > 0 && sl.increasefunc != nil {
		for _, v := range increase[:length] {
			sl.increasefunc(v.stat, v.size)
		}
	}

	if needNotify && sl.addfunc != nil {
		sl.addfunc()
	}
	return
}

func (sl *StatLinked) SetOnAdded(f func()) {
	sl.addfunc = f
}

func (sl *StatLinked) Find(s *Stat) (ns *Stat, ok bool) {
	last, err := sl.Last()
	if err != nil { return }
	if s.SameIno(last) { return last, true }

	for ! s.SameIno(last) {
		last, err = sl.Prev(last)
		if err != nil { return }
	}
	return last, true
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

func (sl *StatLinked) Remove(s *Stat) (ok bool) {
	idx := sl.Search(s)
	if idx < 0 { return }

	sl.statL.Lock()
	defer sl.statL.Unlock()

	sl.stat = append(sl.stat[:idx], sl.stat[idx+1:]...)
	return true
}
