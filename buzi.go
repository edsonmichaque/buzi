package buzi

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"unicode"
)

type Generator interface {
	Apply(map[string]string, *Manifest) ([]File, error)
}

type Metadata struct {
	BaseURL  string `json:"base_url,omitempty" yaml:"base_url,omitempty"`
	Name     string `json:"name,omitempty" yaml:"name,omitempty"`
	FullName string `json:"full_name,omitempty" yaml:"full_name,omitempty"`
	UID      string `json:"uid,omitempty" yaml:"uid,omitempty"`
}

type Manifest struct {
	Metadata   Metadata               `json:"metadata" yaml:"metadata"`
	Operations map[string]Operation   `json:"operations,omitempty" yaml:"operations,omitempty"`
	Types      map[string]*Definition `json:"types,omitempty" yaml:"types,omitempty"`
	Params     map[string]string      `json:"-" yaml:"-"`
}

type Operation struct {
	Name   string        `json:"name" yaml:"name"`
	Http   *Http         `json:"http,omitempty" yaml:"http,omitempty"`
	Input  *Definition   `json:"input,omitempty" yaml:"input,omitempty"`
	Output *Definition   `json:"output,omitempty" yaml:"output,omitempty"`
	Errors []*Definition `json:"errors,omitempty" yaml:"errors,omitempty"`
}

type DefinitionValue struct {
	Type       string                 `json:"type" yaml:"type"`
	Format     string                 `json:"format" yaml:"format"`
	Properties map[string]*Definition `json:"properties,omitempty" yaml:"properties,omitempty"`
	Items      *Definition            `json:"items,omitempty" yaml:"items,omitempty"`
	Error      *Error                 `json:"error,omitempty" yaml:"error,omitempty"`
	Exception  bool                   `json:"exception,omitempty" yaml:"exception,omitempty"`
	Sensitive  bool                   `json:"sensitive,omitempty" yaml:"sensitive,omitempty"`
	Fault      bool                   `json:"fault,omitempty" yaml:"fault,omitempty"`
	Pattern    string                 `json:"pattern,omitempty" yaml:"pattern,omitempty"`
}

type Error struct {
	StatusCode int `json:"status_code" yaml:"status_code"`
}

type Reference struct {
	Ref string `json:"ref" yaml:"ref"`
}

type Definition struct {
	Value *DefinitionValue `json:",inline" yaml:",inline"`
	Ref   *Reference       `json:",inline" yaml:",inline"`
}

type Http struct {
	Method        string        `json:"method" yaml:"method"`
	RequestURI    string        `json:"request_uri" yaml:"request_uri"`
	Authorization Authorization `json:"auth" yaml:"auth"`
}

type Authorization struct {
	Basic *struct {
		User     string
		Password string
	}
	Bearer bool
}

type BasicAuth struct{}

var (
	ErrNotImplemented = errors.New("buzi: not implemented")
)

type File struct {
	Path    string
	Content []byte
}

func (f File) String() string {
	return fmt.Sprintf("Path: %s, Content: %s", f.Path, string(f.Content))
}

func Render(templates fs.FS, tpl string, m *Manifest) ([]byte, error) {
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

func IsObj(k Definition) bool {
	return k.Value != nil && k.Value.Type == "object"
}

func IsScalar(k Definition) bool {
	if k.Value != nil {
		return k.Value.Type != "object" && k.Value.Type != "array"
	}

	return false
}

func IsRef(k Definition) bool {
	return k.Ref != nil
}

func GetRef(k Definition) string {
	return k.Ref.Ref
}

func GetKind(k Definition) string {
	return k.Value.Type
}

func IsArray(k Definition) bool {
	if k.Value != nil {
		return k.Value.Type == "array"
	}

	return false
}

func Props(d Definition) map[string]*Definition {
	props := map[string]*Definition{}
	if d.Value != nil {
		for name, kind := range d.Value.Properties {
			if kind.Value != nil {
				props[name] = kind
			}
		}
	}

	return props
}

func Refs(d Definition) map[string]*Definition {
	props := map[string]*Definition{}
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
