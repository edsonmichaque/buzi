package yaml

import (
	"io"

	"github.com/edsonmichaque/buzi/readers"
	"github.com/edsonmichaque/go-openapi/oas3"
	"gopkg.in/yaml.v3"
)

func New(reader readers.Reader) *parser {
	return &parser{
		reader: reader.GetReader(),
	}
}

type parser struct {
	reader io.Reader
}

func (j *parser) Parse() (*oas3.Spec, error) {
	rawBytes, err := io.ReadAll(j.reader)
	if err != nil {
		return nil, err
	}

	var openapi oas3.Spec

	if err := yaml.Unmarshal(rawBytes, &openapi); err != nil {
		return nil, err
	}

	return &openapi, nil
}
