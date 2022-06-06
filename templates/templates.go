package templates

import (
	"embed"
	"io"
)

//go:embed *
var templates embed.FS

func Get(path string) (io.Reader, error) {
	return templates.Open(path)
}
