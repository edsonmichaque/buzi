package golang

import (
	"path"

	"github.com/edsonmichaque/buzi"
)

type models struct {
	path string
}

func (m models) Apply(params map[string]string, manifest *buzi.Manifest) ([]buzi.File, error) {
	if err := buzi.Require(params, ParamModule, ParamPackage); err != nil {
		return nil, err
	}

	buzi.SetParams(manifest, params)

	f, err := buzi.Render(templates, path.Join("templates", "model.go.tpl"), manifest)
	if err != nil {
		return nil, err
	}

	return []buzi.File{{Path: "models.go", Content: f}}, nil
}
