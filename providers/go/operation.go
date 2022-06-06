package golang

import (
	"errors"

	"github.com/edsonmichaque/umbeluzi"
)

type operationGenerator struct{}

func (o operationGenerator) Generate(spec *umbeluzi.Document) ([]byte, error) {

	return nil, errors.New("buzi: not implemented")
}
