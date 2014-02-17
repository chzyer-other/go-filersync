package filersync

import (
	"io"
)

type Status struct {
	sl *StatLinked
	offset *Offset
	fileList map[uint64] *File
	notifyChange func()
}

func NewStatus(path []string) (s *Status, err error) {
	sl, err := NewStatLinked(path)
	if err != nil { return }
	s = &Status {
		sl: sl,
		offset: NewOffset(),
	}
	sl.SetOnAdded(s.onStatLinkedAdded)
	return
}

func RestoreStatus(r io.Reader) (s *Status, err error) {
	return
}

func (s *Status) Store() {
	return
}

func (s *Status) getFile(stat *Stat) (f *File, offset int64, err error) {
	f, ok := s.fileList[stat.Inode()]
	if ! ok {
		f, err = NewFile(stat)
	}
	offset, ok = s.offset.GetStatOffset(stat)
	if ! ok {
		s.offset.SetStatOffset(stat, 0)
	}
	return
}

func (s *Status) Offset(f *File) (offset int64, err error) {
	stat := f.Stat
	_, offset, err = s.getFile(stat)
	return
}

func (s *Status) Last() (f *File, offset int64, err error) {
	stat, err := s.sl.Last()
	if err != nil { return }
	f, offset, err = s.getFile(stat)
	return
}

func (s *Status) Next(f *File) (nf *File, offset int64, err error) {
	ns, err := s.sl.Next(f.Stat)
	if err != nil { return }
	nf, offset, err = s.getFile(ns)
	return
}

func (s *Status) IsFinish(f *File) (yes bool) {
	return s.offset.IsFinish(f.Stat)
}

func (s *Status) Prev(f *File) (nf *File, offset int64, err error) {
	ps, err := s.sl.Prev(f.Stat)
	if err != nil { return }
	nf, offset, err = s.getFile(ps)
	return
}

func (s *Status) UpdateOffset(f *File, offset int64) {
	s.offset.SetStatOffset(f.Stat, offset)
	if s.IsFinish(f) {
		if s.notifyChange == nil { return }
		s.notifyChange()
	}
	return
}

func (s *Status) UpdatePath(path []string) {
	s.sl.UpdatePath(path)
}

func (s *Status) SetNotifyFunc(notify func()) {
	s.notifyChange = notify
}

func (s *Status) onStatLinkedAdded() {
	if s.notifyChange == nil { return }
	s.notifyChange()
}
