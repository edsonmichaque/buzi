package python

import (
	"path"

	"github.com/edsonmichaque/buzi"
	"github.com/edsonmichaque/buzi/types"
)

const (
	ParamModule = "module"
)

type operations struct{}

func (m operations) Apply(params map[string]string, manifest *types.Manifest) ([]types.File, error) {
	if err := buzi.Require(params, ParamModule); err != nil {
		return nil, err
	}

	buzi.SetParams(manifest, params)

	f, err := buzi.Render(templates, path.Join("templates", "client.py.tpl"), manifest)
	if err != nil {
		return nil, err
	}

	return []types.File{{Path: "client.py", Content: f}}, nil
}
