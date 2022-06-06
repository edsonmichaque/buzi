package python

import (
	"github.com/edsonmichaque/buzi/codegen"
	"github.com/edsonmichaque/buzi/textutil"
)

type provider struct {
	Extension string
}

func New() *provider {
	return &provider{
		Extension: "py",
	}
}

func (r provider) Schema() codegen.ParamsSchema {
	return codegen.ParamsSchema(
		map[string]codegen.ParamSchemaItem{
			"package-name": {
				Type:     "string",
				Required: true,
			},
		},
	)
}

func (r provider) GenerateModels(d codegen.Def, params codegen.OptionSet) ([]codegen.FileDef, error) {
	models := make([]codegen.FileDef, 0)

	for _, model := range d.Models {
		renderModelsParams := codegen.RenderModelParams{
			Template:  "python/model.py.tmpl",
			TargetDir: "lib/model",
			PathFunc:  textutil.ToSnakeCase,
			NameFunc:  textutil.ToPascalCase,
			Extension: r.Extension,
		}

		model, err := codegen.RenderModel(model, renderModelsParams)
		if err != nil {
			return nil, err
		}

		models = append(models, *model)
	}

	return models, nil
}

func (r provider) GenerateOperations(d codegen.Def, params codegen.OptionSet) ([]codegen.FileDef, error) {
	renderOperationsParams := codegen.RenderOperationsParams{
		Template:   "python/client.py.tmpl",
		TargetFile: "client",
		TargetDir:  "lib",
		PathFunc:   textutil.ToSnakeCase,
		NameFunc:   textutil.ToSnakeCase,
		Extension:  r.Extension,
	}

	operations, err := codegen.RenderOperations(d.Operations, renderOperationsParams)
	if err != nil {
		return nil, err
	}

	return operations, nil
}
