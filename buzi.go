package boa

import (
	"github.com/edsonmichaque/buzi/codegen"
	"github.com/edsonmichaque/umbeluzi/types"
)

type Generator interface {
	Generate(def codegen.Def, params codegen.OptionSet) ([]codegen.FileDef, error)
	Schema() codegen.ParamsSchema
}

type Gen interface {
	Generate(types.Spec) ([]byte, error)
}
