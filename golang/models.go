package golang

import (
	"bytes"
	"errors"
	"path"
	"unicode"

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

	if _, ok := params[ParamPackage]; !ok {
		return nil, errors.New("missing param" + ParamPackage)
	}

	mn.Params = params

	f, err := buzi.Render(templates, path.Join("templates", "model.go.tpl"), &mn)
	if err != nil {
		return nil, err
	}

	return []buzi.File{{Path: "models.go", Content: f}}, nil
}

func snakeCase(s string) string {
	buf := new(bytes.Buffer)
	for i, c := range s {
		if unicode.IsUpper(c) {
			if i != 0 && i != len(s)-1 {
				buf.WriteRune('_')
			}
		}

		buf.WriteRune(unicode.ToLower(c))
	}

	return buf.String()
}
