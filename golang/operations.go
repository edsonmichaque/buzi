package golang

import (
	"errors"
	"path"

	"github.com/edsonmichaque/buzi"
)

type operations struct{}

func (m operations) Apply(params map[string]string, manifest *buzi.Manifest) ([]buzi.File, error) {
	mn := *manifest

	if _, ok := params[ParamModule]; !ok {
		return nil, errors.New("missing param" + ParamModule)
	}

	if _, ok := params[ParamPackage]; !ok {
		return nil, errors.New("missing param" + ParamPackage)
	}

	mn.Params = params

	f, err := buzi.Render(templates, path.Join("templates", "client.go.tpl"), &mn)
	if err != nil {
		return nil, err
	}

	return []buzi.File{{Path: "client.go", Content: f}}, nil
}
