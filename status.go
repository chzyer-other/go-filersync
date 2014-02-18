package filersync

import (
	"io"
)

type Status struct {
	sl *StatLinked
	oset *Offset
	fileList map[uint64] *File
	notifyChange func()
}

func NewStatus(path []string) (s *Status, err error) {
	sl, err := NewStatLinked(path)
	if err != nil { return }
	s = &Status {
		sl: sl,
		oset: NewOffset(),
	}
	sl.SetOnAdded(s.onStatLinkedAdded)
	// sl.SetOnSizeIncrease(s.onStatLinkedSizeIncrease)
	return
}

func RestoreStatus(r io.Reader) (s *Status, err error) {
	return
}

func (s *Status) Store() {
	return
}

func (s *Status) getFile(stat *Stat) (f *File, err error) {
	f, ok := s.fileList[stat.Inode()]
	if ! ok {
		f, err = NewFile(stat)
	}
	return
}

func (s *Status) Offset(f *File) (offset int64) {
	return s.offset(f.Stat)
}

func (s *Status) offset(stat *Stat) (offset int64) {
	offset, ok := s.oset.GetStatOffset(stat)
	if ! ok {
		s.oset.SetStatOffset(stat, 0)
	}
	return
}

func (s *Status) Last() (f *File, err error) {
	stat, err := s.sl.Last()
	if err != nil { return }
	f, err = s.getFile(stat)
	return
}

func (s *Status) Next(f *File) (nf *File, err error) {
	ns, err := s.sl.Next(f.Stat)
	if err != nil { return }
	nf, err = s.getFile(ns)
	return
}

func (s *Status) IsFinish(f *File) (yes bool) {
	s.Sync(f)
	return s.isFinish(f.Stat)
}

func (s *Status) isFinish(ns *Stat) (yes bool) {
	ts, ok := s.sl.Find(ns)
	if ok {
		ns.Size = ts.Size
	}
	offset := s.offset(ns)
	return offset >= ns.Size
}

func (s *Status) Sync(f *File) (ok bool) {
	stat, ok := s.sl.Find(f.Stat)
	if ! ok { return }
	f.Stat = stat
	return true
}

func (s *Status) Prev(f *File) (nf *File, err error) {
	ps, err := s.sl.Prev(f.Stat)
	if err != nil { return }
	nf, err = s.getFile(ps)
	return
}

func (s *Status) UpdateOffset(f *File, offset int64) {
	s.updateStatOffset(f.Stat, offset)
}

func (s *Status) updateStatOffset(ns *Stat, offset int64) {
	s.oset.SetStatOffset(ns, offset)
	if s.isFinish(ns) {
		if s.notifyChange == nil { return }
		s.notifyChange()
	}
}

func (s *Status) AddOffset(f *File, offset int64) {
	s.addOffset(f.Stat, offset)
}

func (s *Status) addOffset(ns *Stat, offset int64) {
	s.oset.AddStatOffset(ns, offset)
	if s.isFinish(ns) {
		if s.notifyChange == nil { return }
		s.notifyChange()
	}
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
