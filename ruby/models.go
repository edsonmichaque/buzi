package ruby

import (
	"errors"
	"path"

	"github.com/edsonmichaque/buzi"
)

type models struct {
	path string
}

func (m models) Apply(params map[string]string, manifest *buzi.Manifest) ([]buzi.File, error) {
	mn := *manifest

	if _, ok := params[ParamModule]; !ok {
		return nil, errors.New("missing param" + ParamModule)
	}

	mn.Params = params

	f, err := buzi.Render(templates, path.Join("templates", "model.rb.tpl"), &mn)
	if err != nil {
		return nil, err
	}

	return []buzi.File{{Path: "models.rb", Content: f}}, nil
}
