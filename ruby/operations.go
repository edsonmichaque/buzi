package ruby

import (
	"errors"
	"path"

	"github.com/edsonmichaque/buzi"
)

const (
	ParamModule = "module"
)

type operations struct{}

func (m operations) Apply(params map[string]string, manifest *buzi.Manifest) ([]buzi.File, error) {
	mn := *manifest

	if _, ok := params[ParamModule]; !ok {
		return nil, errors.New("missing param" + ParamModule)
	}

	mn.Params = params

	f, err := buzi.Render(templates, path.Join("templates", "client.rb.tpl"), &mn)
	if err != nil {
		return nil, err
	}

	return []buzi.File{{Path: "client.rb", Content: f}}, nil
}
