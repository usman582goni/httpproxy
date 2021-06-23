package stream

import (
	"io"
	"sync"
)

type Buffer struct {
	buf []byte
	rw  *sync.RWMutex
}

type bufferReader struct {
	buffer *Buffer
	i      int64
}

func NewBuffer() Source {
	return &Buffer{
		buf: make([]byte, 0),
		rw:  new(sync.RWMutex),
	}
}

func (b *Buffer) Open() (io.Reader, error) {
	return &bufferReader{
		buffer: b,
		i:      0,
	}, nil
}
func (b *bufferReader) Read(dst []byte) (n int, err error) {
	b.buffer.rw.RLock()
	defer b.buffer.rw.RUnlock()
	if b.i >= int64(len(b.buffer.buf)) {
		
		return 0, io.EOF
	}
	n = copy(dst, b.buffer.buf[b.i:])
	b.i += int64(n)
	return n, nil
}

func (b *Buffer) Write(src []byte) (n int, err error) {
	b.rw.Lock()
	defer b.rw.Unlock()
	b.buf = append(b.buf, src...)
	return len(src), nil
}