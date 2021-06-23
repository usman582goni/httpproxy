package stream

import (
	"io"
	"os"
)

type File struct {
	*os.File
}

func NewFile(path string) Source {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	return &File{file}
}

func (f File) Open() (io.Reader, error) {
	return os.Open(f.Name())
}