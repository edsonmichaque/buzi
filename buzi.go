package boa

import (
	"github.com/edsonmichaque/buzi/codegen"
)

type Generator interface {
	Generate(def codegen.Def, params codegen.OptionSet) ([]codegen.FileDef, error)
	Schema() codegen.ParamsSchema
}

// type Renderer interface {
// 	Render(codegen.Def, []codegen.Option) (map[string][]byte, error)
// }

// type ConfigGenerator interface {
// 	GenerateConfig(*oas3.Spec, []codegen.Option) (*codegen.ConfigDef, error)
// }

// type ModelsGenerator interface {
// 	GenerateModels(*oas3.Spec, []codegen.Option) ([]codegen.ModelDef, error)
// }

// type OperationsGenerator interface {
// 	GenerateOperations(*oas3.Spec, []codegen.Option) ([]codegen.OperationDef, error)
// }

// type TypeMapper interface {
// 	Map(string, string) (*string, error)
// }
