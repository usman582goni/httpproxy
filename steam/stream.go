package stream

import (
	"errors"
	"io"
	"sync"
)

func New(source Source) *Stream {
	return &Stream{
		Source: source,
		writer: source,
		cond:   sync.NewCond(new(sync.Mutex)),
		closed: false,
	}
}

type Source interface {
	io.Writer
	Open() (io.Reader, error)
}

type Stream struct {
	Source
	writer io.Writer
	cond   *sync.Cond
	closed bool
}

type streamReader struct {
	stream *Stream
	reader io.Reader
}

func (s *Stream) NewReader() (io.Reader, error) {
	r, err := s.Open()

	if err != nil {
		return nil, err
	}

	return &streamReader{
		stream: s,
		reader: r,
	}, nil
}

func (s *streamReader) Read(p []byte) (n int, err error) {
	s.stream.cond.L.Lock()
	defer s.stream.cond.L.Unlock()
	n, err = s.reader.Read(p)

	
	if err == nil {
		return
	}

	
	for err == io.EOF {
		
		if s.stream.closed {
			return
		}

		if n > 0 {
			return n, nil
		}

		s.stream.cond.Wait()
		n, err = s.reader.Read(p)
	}
	return
}

func (s *Stream) Write(data []byte) (int, error) {
	defer s.cond.Broadcast()
	return s.writer.Write(data)
}

func (s *Stream) CloseWrite() error {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()
	defer s.cond.Broadcast()
	if s.closed {
		return errors.New("Multireader closed multiple times")
	}
	s.closed = true
	return nil
}
