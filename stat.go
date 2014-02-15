package filersync

import (
	"time"
	"sort"
	"syscall"
)

type Stat struct {
	*syscall.Stat_t
	Ctime time.Time
	Mtime time.Time
	Path string
}

func (s *Stat) MBefore(ns *Stat) bool {
	return s.Mtime.Before(ns.Mtime)
}

func (s *Stat) MAfter(ns *Stat) bool {
	return s.Mtime.After(ns.Mtime)
}

func (s *Stat) Inode() uint64 { return s.Ino }

type StatSlice []*Stat

func (p StatSlice) Len() int { return len(p) }
func (p StatSlice) Less(i, j int) bool { return p[i].Mtime.Before(p[j].Mtime) }
func (p StatSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Sort is a convenience method.
func (p StatSlice) Sort() { sort.Sort(p) }

func fstat(fd uintptr) (stat *Stat, err error) {
	var s_stat syscall.Stat_t
	err = syscall.Fstat(int(fd), &s_stat)

	stat = statEncode(&s_stat)
	return
}

func getStat(name string) (stat *Stat, err error) {
	var s_stat syscall.Stat_t
	err = syscall.Stat(name, &s_stat)
	stat = statEncode(&s_stat)
	stat.Path = name
	return
}
