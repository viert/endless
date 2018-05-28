package endless

import (
	"errors"
	"sync"
)

type (
	// Endless is an Endless (ring) buffer
	Endless struct {
		sync.RWMutex
		data        []byte
		writeCursor uint64
		start       uint64
	}

	// Reader represents a Reader interface for an Endless buffer
	Reader struct {
		source *Endless
		pos    uint64
	}
)

// NewEndless creates a ring buffer of a given _size_
func NewEndless(size int) *Endless {
	return &Endless{
		data:        make([]byte, size),
		writeCursor: 0,
		start:       0,
	}
}

func (e *Endless) pos() int {
	return int(e.writeCursor % uint64(len(e.data)))
}

// MidPoint returns the midpoint of the buffer
func (e *Endless) MidPoint() uint64 {
	return (e.start + e.writeCursor) / 2
}

// Start returns the index of the first element in buffer
func (e *Endless) Start() uint64 {
	return e.start
}

// End returns the index of the last element in buffer + 1
// (index of the next writable element)
func (e *Endless) End() uint64 {
	return e.writeCursor
}

// Filled returns true, if the whole buffer contains useful data
// Filled() is false only for a couple of writes in the very beginning
// of the buffer lifecycle.
func (e *Endless) Filled() bool {
	return e.writeCursor-e.start == uint64(len(e.data))
}

func (e *Endless) Write(buf []byte) (int, error) {
	e.Lock()
	defer e.Unlock()
	maxSize := len(buf)
	if maxSize > len(e.data) {
		maxSize = len(e.data)
	}

	if e.pos()+maxSize > len(e.data) {
		tailSize := len(e.data) - e.pos()
		headSize := maxSize - tailSize
		copy(e.data[e.pos():], buf[:tailSize])
		copy(e.data[:headSize], buf[tailSize:])
	} else {
		copy(e.data[e.pos():e.pos()+maxSize], buf[:maxSize])
	}

	e.writeCursor += uint64(maxSize)
	if e.writeCursor-e.start > uint64(len(e.data)) {
		e.start = e.writeCursor - uint64(len(e.data))
	}

	return maxSize, nil
}

// NewReader creates a Reader interface for a given Endless buffer
func (e *Endless) NewReader(start uint64) *Reader {
	reader := new(Reader)
	reader.source = e
	if e.start > start {
		reader.pos = e.start
	} else {
		reader.pos = start
	}
	return reader
}

func (r *Reader) Read(buf []byte) (n int, err error) {
	r.source.RLock()
	defer r.source.RUnlock()
	if r.source.start > r.pos {
		return 0, errors.New("reader has remained behind the buffer")
	}

	maxSize := len(buf)
	if maxSize > int(r.source.writeCursor-r.pos) {
		maxSize = int(r.source.writeCursor - r.pos)
	}

	bpos := int(r.pos % uint64(len(r.source.data)))
	if bpos+maxSize > len(r.source.data) {
		tailSize := len(r.source.data) - bpos
		headSize := maxSize - tailSize
		copy(buf[:tailSize], r.source.data[bpos:])
		copy(buf[tailSize:], r.source.data[:headSize])
	} else {
		copy(buf, r.source.data[bpos:bpos+maxSize])
	}
	r.pos += uint64(maxSize)

	return maxSize, nil
}
