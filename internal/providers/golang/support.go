package golang

import (
	"errors"
	"path/filepath"

	"github.com/edsonmichaque/buzi"
	"github.com/edsonmichaque/buzi/types"
)

type support struct {
	makefileTemplate     string
	gosumTemplate        string
	gomodTemplate        string
	readmeTemplate       string
	contributingTemplate string
}

func (m support) Apply(params map[string]string, manifest *types.Manifest) ([]types.File, error) {
	mn := *manifest

	if _, ok := params[ParamModule]; !ok {
		return nil, errors.New("missing param" + ParamModule)
	}

	if _, ok := params[ParamPackage]; !ok {
		return nil, errors.New("missing param" + ParamPackage)
	}
	mn.Params = params

	gomod, err := buzi.Render(templates, filepath.Join("templates", "go.mod.tpl"), &mn)
	if err != nil {
		return nil, err
	}

	files := []types.File{
		{
			Path:    "go.mod",
			Content: gomod,
		},
	}

	return files, nil
}
