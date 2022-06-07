package ruby

import (
	"path"

	"github.com/edsonmichaque/buzi"
	"github.com/edsonmichaque/buzi/types"
)

type models struct {
	path string
}

func (m models) Apply(params map[string]string, manifest *types.Manifest) ([]types.File, error) {
	if err := buzi.Require(params, ParamModule); err != nil {
		return nil, err
	}

	buzi.SetParams(manifest, params)

	f, err := buzi.Render(templates, path.Join("templates", "model.rb.tpl"), manifest)
	if err != nil {
		return nil, err
	}

	return []types.File{{Path: "models.rb", Content: f}}, nil
}
