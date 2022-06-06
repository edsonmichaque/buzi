package codegen

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/edsonmichaque/buzi/templates"
	"github.com/edsonmichaque/go-openapi/oas3"
)

type Def struct {
	Config     ConfigDef
	Models     []ModelDef
	Operations []OperationDef
}

type ConfigDef struct {
	Name       string
	Properties []PropertyDef
}

type ModelDef struct {
	Name        string
	Description string
	Type        string
	Properties  []PropertyDef
}

type PropertyDef struct {
	Name        string
	Type        string
	Annotations map[string]interface{}
}

type OperationDef struct {
	Name        string
	Description string
	Summary     string
	Method      string
	Path        string
	Params      []oas3.Parameter
}

type Option struct {
	Name  string
	Value interface{}
}

func NewOperationsDef(spec *oas3.Spec, o []Option) ([]OperationDef, error) {
	operationsDef := make([]OperationDef, 0)

	for path, pathItem := range spec.Paths {
		for method, operation := range pathItem.Operations() {
			newOperation := OperationDef{
				Name:        operation.OperationID,
				Summary:     operation.Summary,
				Description: operation.Description,
				Method:      method,
				Path:        path,
				Params:      operation.Parameters,
			}

			operationsDef = append(operationsDef, newOperation)
		}
	}

	return operationsDef, nil
}

func NewModelsDef(spec *oas3.Spec, o []Option) ([]ModelDef, error) {
	models := make([]ModelDef, 0)

	for name, _ := range spec.Components.Schemas {
		model := ModelDef{
			Name: name,
		}

		models = append(models, model)
	}

	return models, nil
}

func NewDef(spec *oas3.Spec) (*Def, error) {
	operationsDef, err := NewOperationsDef(spec, nil)
	if err != nil {
		return nil, err
	}

	modelsDef, err := NewModelsDef(spec, nil)
	if err != nil {
		return nil, err
	}

	defs := &Def{
		Operations: operationsDef,
		Models:     modelsDef,
	}

	return defs, nil
}

type FileDef struct {
	Path    string
	Content []byte
}

type RenderModelParams struct {
	Template  string
	TargetDir string
	PathFunc  func(string) string
	NameFunc  func(string) string
	Extension string
}

func RenderModel(m ModelDef, params RenderModelParams) (*FileDef, error) {
	tpl, err := templates.Get(params.Template)
	if err != nil {
		return nil, err
	}

	rawBytes, err := io.ReadAll(tpl)
	if err != nil {
		return nil, err
	}

	newTpl, err := template.New("model").Parse(string(rawBytes))
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	modelDefCopy := ModelDef{
		Name:        params.NameFunc(m.Name),
		Type:        m.Type,
		Description: m.Description,
	}

	if err := newTpl.Execute(buf, modelDefCopy); err != nil {
		return nil, err
	}

	var path string
	if params.TargetDir != "" {
		path = filepath.Join(params.TargetDir, fmt.Sprintf("%s/%s.%s", params.TargetDir, params.PathFunc(m.Name), params.Extension))
	} else {
		path = fmt.Sprintf("%s/%s.%s", params.TargetDir, params.PathFunc(m.Name), params.Extension)
	}

	file := &FileDef{
		Path:    path,
		Content: buf.Bytes(),
	}

	return file, nil
}

type RenderOperationsParams struct {
	Template   string
	TargetDir  string
	TargetFile string
	PathFunc   func(string) string
	NameFunc   func(string) string
	Extension  string
}

func RenderOperations(ops []OperationDef, params RenderOperationsParams) ([]FileDef, error) {
	tpl, err := templates.Get(params.Template)
	if err != nil {
		return nil, err
	}

	rawBytes, err := io.ReadAll(tpl)
	if err != nil {
		return nil, err
	}

	newTpl, err := template.New("client").Parse(string(rawBytes))
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	newOps := make([]OperationDef, 0)
	for _, op := range ops {
		eachOp := OperationDef{
			Name:        params.NameFunc(op.Name),
			Path:        op.Path,
			Method:      op.Method,
			Summary:     op.Summary,
			Description: op.Description,
			Params:      op.Params,
		}

		newOps = append(newOps, eachOp)
	}

	if err := newTpl.Execute(buf, newOps); err != nil {
		return nil, err
	}

	var path string
	if params.TargetDir != "" {
		path = filepath.Join(params.TargetDir, params.PathFunc(params.TargetFile)+"."+params.Extension)
	} else {
		path = params.PathFunc(params.TargetFile) + "." + params.Extension
	}

	file := FileDef{
		Path:    path,
		Content: buf.Bytes(),
	}

	return []FileDef{file}, nil
}

type RenderParams struct {
	Template   string
	TargetDir  string
	TargetFile string
	PathFunc   func(string) string
	NameFunc   func(string) string
	Extension  string
}

func RenderExtra(data map[string]interface{}, params RenderParams) (*FileDef, error) {
	tpl, err := templates.Get(params.Template)
	if err != nil {
		return nil, err
	}

	rawBytes, err := io.ReadAll(tpl)
	if err != nil {
		return nil, err
	}

	newTpl, err := template.New("extra").Parse(string(rawBytes))
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	if err := newTpl.Execute(buf, data); err != nil {
		return nil, err
	}

	var path string
	if params.Extension != "" {
		path = params.TargetFile + "." + params.Extension
	} else {
		path = params.TargetFile
	}

	file := &FileDef{
		Path:    path,
		Content: buf.Bytes(),
	}

	return file, nil
}

type ParamSchemaItem struct {
	Type     string
	Required bool
}

type ParamsSchema map[string]ParamSchemaItem

type OptionSet map[string]interface{}

func OptionSetFromMap(m map[string]interface{}) OptionSet {
	return OptionSet(m)
}

func (o OptionSet) Get(key string) (interface{}, bool) {
	m := map[string]interface{}(o)

	if v, ok := m[key]; ok {
		return v, ok
	}

	return nil, false
}

func (o OptionSet) GetString(key string) (string, bool) {
	if v, ok := o.Get(key); ok {
		if newStr, ok := v.(string); ok {
			return newStr, true
		}
	}

	return "", false
}

func (o OptionSet) GetInt(key string) (int, bool) {
	if v, ok := o.Get(key); ok {
		if newInt, ok := v.(int); ok {
			return newInt, true
		}
	}

	return 0, false
}

func (o OptionSet) GetBool(key string) (bool, bool) {
	if v, ok := o.Get(key); ok {
		if newBool, ok := v.(bool); ok {
			return newBool, true
		}
	}

	return false, false
}

func (o OptionSet) GetFloat64(key string) (float64, bool) {
	if v, ok := o.Get(key); ok {
		if newFloat64, ok := v.(float64); ok {
			return newFloat64, true
		}
	}

	return 0, false
}

func (o OptionSet) GetFloat32(key string) (float32, bool) {
	if v, ok := o.Get(key); ok {
		if newFloat32, ok := v.(float32); ok {
			return newFloat32, true
		}
	}

	return 0, false
}

type ValidationError struct {
	errs []string
}

func (v ValidationError) Error() string {
	return strings.Join(v.errs, "\n")
}

func ValidateParams(m map[string]interface{}, schema ParamsSchema) error {
	ss := map[string]ParamSchemaItem(schema)

	errs := make([]string, 0)
	for k, v := range ss {
		if kk, ok := m[k]; ok {
			switch v.Type {
			case "string":
				if _, ok := kk.(string); !ok {
					errs = append(errs, k+" is required")
				}
			case "int":
				if _, ok := kk.(int); !ok {
					errs = append(errs, k+" is required")
				}
			case "int32":
				if _, ok := kk.(int32); !ok {
					errs = append(errs, k+" is required")
				}
			case "int64":
				if _, ok := kk.(int64); !ok {
					errs = append(errs, k+" is required")
				}
			case "float32":
				if _, ok := kk.(float32); !ok {
					errs = append(errs, k+" is required")
				}
			case "float64":
				if _, ok := kk.(float64); !ok {
					errs = append(errs, k+" is required")
				}
			case "bool":
				if _, ok := kk.(bool); !ok {
					errs = append(errs, k+" is required")
				}
			}

		} else {
			if v.Required {
				errs = append(errs, k+" is required")
			}
		}
	}

	if len(errs) > 0 {
		return ValidationError{
			errs: errs,
		}
	}

	return nil
}

func Generate(p Generator, def Def, params map[string]interface{}) ([]FileDef, error) {
	if err := ValidateParams(params, p.Schema()); err != nil {
		return nil, err
	}

	// config, err := p.GenerateConfig(def, params)
	// if err != nil {
	// 	return nil, err
	// }

	models, err := p.GenerateModels(def, params)
	if err != nil {
		return nil, err
	}

	operations, err := p.GenerateOperations(def, params)
	if err != nil {
		return nil, err
	}

	return MergeFileDefs( /*config,*/ models, operations), nil
}

func MergeFileDefs(defs ...[]FileDef) []FileDef {
	l := make([]FileDef, 0)

	for _, d := range defs {
		l = append(l, d...)
	}

	return l
}

type Generator interface {
	// GenerateConfig(def Def, params OptionSet) ([]FileDef, error)
	GenerateModels(def Def, params OptionSet) ([]FileDef, error)
	GenerateOperations(def Def, params OptionSet) ([]FileDef, error)
	Schema() ParamsSchema
}
