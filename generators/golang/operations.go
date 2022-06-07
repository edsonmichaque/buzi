package golang

import (
	"path"

	"github.com/edsonmichaque/buzi"
	"github.com/edsonmichaque/buzi/types"
)

type operations struct{}

func (m operations) Apply(params map[string]string, manifest *types.Manifest) ([]types.File, error) {
	if err := buzi.Require(params, ParamModule, ParamPackage); err != nil {
		return nil, err
	}

	buzi.SetParams(manifest, params)

	f, err := buzi.Render(templates, path.Join("templates", "client.go.tpl"), manifest)
	if err != nil {
		return nil, err
	}

	return []types.File{{Path: "client.go", Content: f}}, nil
}
