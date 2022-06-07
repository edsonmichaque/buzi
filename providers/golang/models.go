package golang

import (
	"path"

	"github.com/edsonmichaque/buzi"
	"github.com/edsonmichaque/buzi/types"
)

type models struct {
	path string
}

func (m models) Apply(params map[string]string, manifest *types.Manifest) ([]types.File, error) {
	if err := buzi.Require(params, ParamModule, ParamPackage); err != nil {
		return nil, err
	}

	buzi.SetParams(manifest, params)

	f, err := buzi.Render(templates, path.Join("templates", "model.go.tpl"), manifest)
	if err != nil {
		return nil, err
	}

	return []types.File{{Path: "models.go", Content: f}}, nil
}
