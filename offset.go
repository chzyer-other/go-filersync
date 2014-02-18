package filersync

import (
	"sync"
)

type statData struct {
	offset int64
	finish bool
	limit int64
}

func (s *statData) SetOffset(offset, size int64) {
	s.offset = offset
	s.finish = s.offset >= size
}

func (s *statData) SetLimit(limit int64) {
	s.limit = limit
	s.finish = s.offset >= limit
}

type Offset struct {
	data map[uint64] *statData // map[ino] offset
	l sync.Mutex
}

func NewOffset() (o *Offset) {
	return &Offset {
		data: make(map[uint64] *statData),
	}
}

func (o *Offset) GetStatOffset(f *Stat) (offset int64, ok bool) {
	ino := f.Inode()
	o.l.Lock()
	defer o.l.Unlock()
	return o.getOffset(ino)
}

func (o *Offset) getOffset(ino uint64) (offset int64, ok bool) {
	sd, ok := o.data[ino]
	if ! ok { return }
	offset = sd.offset
	return
}

func (o *Offset) SetStatOffset(f *Stat, offset int64) {
	ino := f.Inode()

	o.l.Lock()
	defer o.l.Unlock()
	if offset > f.Size { offset = f.Size }
	o.setOffset(ino, offset, f.Size)
}

func (o *Offset) AddStatOffset(f *Stat, offset int64) {
	ino := f.Inode()
	o.l.Lock()
	defer o.l.Unlock()
	d, ok := o.data[ino]
	if ! ok {
		o.setOffset(ino, offset, f.Size)
		return
	}
	d.SetOffset(d.offset + offset, f.Size)
}

func (o *Offset) setOffset(ino uint64, offset, size int64) {
	_, ok := o.data[ino]
	if ok {
		o.data[ino].SetOffset(offset, size)
	} else {
		o.data[ino] = &statData {
			offset: offset,
			finish: offset >= size,
		}
	}
}

func (o *Offset) SetStatLimit(f *Stat, limit int64) (ok bool) {
	ino := f.Inode()
	o.l.Lock()
	defer o.l.Unlock()
	return o.setStatLimit(ino, limit, f.Size)
}

func (o *Offset) setStatLimit(ino uint64, limit, size int64) (ok bool) {
	d, ok := o.data[ino]
	if ! ok { return }
	d.SetLimit(limit)
	return
}

func (o *Offset) IsFinish(f *Stat) (ok bool) {
	ino := f.Inode()
	o.l.Lock()
	defer o.l.Unlock()
	return o.isFinish(ino)
}

func (o *Offset) isFinish(ino uint64) (ok bool) {
	sd, ok := o.data[ino]
	if ! ok { return }
	return sd.finish
}

func (o *Offset) Remove(f *Stat) (ok bool) {
	ino := f.Inode()
	o.l.Lock()
	defer o.l.Unlock()
	return o.remove(ino)
}

func (o *Offset) remove(ino uint64) (ok bool) {
	_, ok = o.data[ino]
	if ! ok { return }
	delete(o.data, ino)
	return
}
