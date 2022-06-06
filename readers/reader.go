package readers

import "io"

type Reader interface {
	GetReader() io.Reader
}
