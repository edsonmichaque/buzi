package local

import (
	"io"
	"os"
)

type local struct {
	reader io.Reader
}

func New(path string) (*local, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return &local{reader: file}, nil
}

func (f *local) GetReader() io.Reader {
	return f.reader
}
