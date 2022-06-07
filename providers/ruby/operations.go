package ruby

import (
	"path"

	"github.com/edsonmichaque/buzi"
)

const (
	ParamModule = "module"
)

type operations struct{}

func (m operations) Apply(params map[string]string, manifest *buzi.Manifest) ([]buzi.File, error) {
	if err := buzi.Require(params, ParamModule); err != nil {
		return nil, err
	}

	buzi.SetParams(manifest, params)

	f, err := buzi.Render(templates, path.Join("templates", "client.rb.tpl"), manifest)
	if err != nil {
		return nil, err
	}

	return []buzi.File{{Path: "client.rb", Content: f}}, nil
}
