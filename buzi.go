package buzi

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"io/fs"
	"unicode"

	"github.com/edsonmichaque/buzi/types"
)

type Generator interface {
	Apply(map[string]string, *types.Manifest) ([]types.File, error)
}

func Render(templates fs.FS, tpl string, m *types.Manifest) ([]byte, error) {
	f, err := templates.Open(tpl)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, f); err != nil {
		return nil, err
	}

	t, err := template.New("model").Funcs(FuncMap).Parse(buf.String())
	if err != nil {
		return nil, err
	}

	buf2 := new(bytes.Buffer)
	if err := t.Execute(buf2, m); err != nil {
		return nil, err
	}

	return buf2.Bytes(), nil
}

var FuncMap = template.FuncMap{
	"camelcase": camelCase,
	"isref":     IsRef,
	"isobj":     IsObj,
	"isscalar":  IsScalar,
	"props":     Props,
	"isarray":   IsArray,
	"kind":      GetKind,
	"ref":       GetRef,
	"refs":      Refs,
	"snakecase": snakeCase,
	"titlecase": titleCase,
}

func IsObj(k types.Definition) bool {
	return k.Value != nil && k.Value.Type == "object"
}

func IsScalar(k types.Definition) bool {
	if k.Value != nil {
		return k.Value.Type != "object" && k.Value.Type != "array"
	}

	return false
}

func IsRef(k types.Definition) bool {
	return k.Ref != nil
}

func GetRef(k types.Definition) string {
	return k.Ref.Ref
}

func GetKind(k types.Definition) string {
	return k.Value.Type
}

func IsArray(k types.Definition) bool {
	if k.Value != nil {
		return k.Value.Type == "array"
	}

	return false
}

func Props(d types.Definition) map[string]*types.Definition {
	props := map[string]*types.Definition{}
	if d.Value != nil {
		for name, kind := range d.Value.Properties {
			if kind.Value != nil {
				props[name] = kind
			}
		}
	}

	return props
}

func Refs(d types.Definition) map[string]*types.Definition {
	props := map[string]*types.Definition{}
	if d.Value != nil {
		for name, kind := range d.Value.Properties {
			if kind.Ref != nil {
				props[name] = kind
			}
		}
	}

	return props
}

func camelCase(s string) string {
	buf := new(bytes.Buffer)
	for i, c := range s {

		if unicode.IsUpper(c) {
			if i != 0 && i != len(s)-1 {
				buf.WriteRune('_')
			}
		}

		if i == 0 {
			buf.WriteRune(unicode.ToUpper(c))
		} else {
			buf.WriteRune(c)
		}
	}

	return buf.String()
}

func titleCase(s string) string {
	buf := new(bytes.Buffer)
	for i, c := range s {

		if unicode.IsUpper(c) {
			if i != 0 && i != len(s)-1 {
				buf.WriteRune('_')
			}
		}

		if i == 0 {
			buf.WriteRune(unicode.ToUpper(c))
		} else {
			buf.WriteRune(c)
		}
	}

	return buf.String()
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

func Require(params map[string]string, list ...string) error {
	for _, s := range list {
		if _, ok := params[s]; !ok {
			return errors.New("buzi: missing param " + s)
		}
	}

	return nil
}

func SetParams(mn *types.Manifest, p map[string]string) {
	if mn.Params == nil {
		mn.Params = p
	}
}
