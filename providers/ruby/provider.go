package ruby

import (
	"github.com/edsonmichaque/buzi/codegen"
	"github.com/edsonmichaque/buzi/textutil"
)

type provider struct{}

func New() *provider {
	return &provider{}
}

func (r provider) Schema() codegen.ParamsSchema {
	return codegen.ParamsSchema(
		map[string]codegen.ParamSchemaItem{
			"gem-name": {
				Type:     "string",
				Required: true,
			},
			"module-name": {
				Type:     "string",
				Required: true,
			},
		},
	)
}

func (r provider) GenerateModels(def codegen.Def, params codegen.OptionSet) ([]codegen.FileDef, error) {
	models := make([]codegen.FileDef, 0)

	for _, model := range def.Models {
		renderModelsParams := codegen.RenderModelParams{
			Template:  "ruby/model.rb.tmpl",
			TargetDir: "lib/model",
			PathFunc:  textutil.SnakeCase,
			NameFunc:  textutil.ToPascalCase,
			Extension: "rb",
		}

		model, err := codegen.RenderModel(model, renderModelsParams)
		if err != nil {
			return nil, err
		}

		models = append(models, *model)
	}

	return models, nil
}

func (r provider) GenerateOperations(def codegen.Def, params codegen.OptionSet) ([]codegen.FileDef, error) {
	renderOperationsParams := codegen.RenderOperationsParams{
		Template:   "ruby/client.rb.tmpl",
		TargetFile: "client",
		TargetDir:  "lib",
		PathFunc:   textutil.SnakeCase,
		NameFunc:   textutil.SnakeCase,
		Extension:  "rb",
	}

	operations, err := codegen.RenderOperations(def.Operations, renderOperationsParams)
	if err != nil {
		return nil, err
	}

	return operations, nil
}
