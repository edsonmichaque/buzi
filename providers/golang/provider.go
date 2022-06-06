package golang

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"

	boa "github.com/edsonmichaque/buzi"
	"github.com/edsonmichaque/buzi/stringutil"
	"github.com/edsonmichaque/buzi/templates"
	"github.com/edsonmichaque/go-openapi/oas3"
)

type provider struct{}

func New() *provider {
	return &provider{}
}

func (p provider) Generate(spec *oas3.Spec) (*boa.SDK, error) {
	models, err := p.GenerateModels(spec)
	if err != nil {
		return nil, err
	}

	client, err := p.GenerateOperations(spec)
	if err != nil {
		return nil, err
	}

	for k, v := range client {
		models[k] = v
	}

	return &boa.SDK{Files: models}, nil
}

func (p provider) GenerateModels(spec *oas3.Spec) (map[string][]byte, error) {

	models := make(map[string][]byte)

	for name, schema := range spec.Components.Schemas {
		params := map[string]string{
			"Name":        name,
			"Description": schema.Description,
		}

		tpl, err := templates.Get("golang/model.go.tmpl")
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

		if err := newTpl.Execute(buf, params); err != nil {
			return nil, err
		}

		models[fmt.Sprintf("go/models/%s.go", boa.ToSnakeCase(name))] = buf.Bytes()
	}

	return models, nil
}

func (p provider) GenerateOperations(spec *oas3.Spec) (map[string][]byte, error) {

	files := make(map[string][]byte)

	params := make([]operation, 0)

	for path, pathItem := range spec.Paths {

		for method, op := range pathItem.Operations() {
			newOp := operation{
				Name:        stringutil.ToPascalCase(op.OperationID),
				Description: op.Description,
				Method:      method,
				Path:        path,
				PathParams:  make([]param, 0),
			}

			for _, pr := range op.Parameters {
				if pr.In == "path" {

					newParam := param{
						Name:         stringutil.ToPascalCase(pr.Name),
						Type:         "path",
						OriginalName: stringutil.ToCamelCase(pr.Name),
						Template:     "{" + pr.Name + "}",
						Schema:       *pr.Schema,
					}

					newOp.PathParams = append(newOp.PathParams, newParam)
				} else if pr.In == "query" {
					newParam := param{
						Name:         stringutil.ToPascalCase(pr.Name),
						Type:         "query",
						OriginalName: pr.Name,
						Template:     "{" + pr.Name + "}",
						Schema:       *pr.Schema,
					}

					newOp.QueryParams = append(newOp.QueryParams, newParam)
				} else if pr.In == "header" {
					newParam := param{
						Name:         stringutil.ToPascalCase(pr.Name),
						Type:         "header",
						OriginalName: pr.Name,
						Template:     "{" + pr.Name + "}",
						Schema:       *pr.Schema,
					}

					newOp.HeaderParams = append(newOp.HeaderParams, newParam)
				} else if pr.In == "cookie" {
					newParam := param{
						Name:         stringutil.ToPascalCase(pr.Name),
						Type:         "cookie",
						OriginalName: pr.Name,
						Template:     "{" + pr.Name + "}",
						Schema:       *pr.Schema,
					}

					newOp.CookieParams = append(newOp.CookieParams, newParam)
				}
			}

			params = append(params, newOp)
		}
	}

	tpl, err := templates.Get("golang/client.go.tmpl")
	if err != nil {
		return nil, err
	}

	rawBytes, err := io.ReadAll(tpl)
	if err != nil {
		return nil, err
	}

	newTpl, err := template.New("client").Funcs(template.FuncMap{
		"detectType":      DetectType,
		"convertToString": ConvertToString,
	}).Parse(string(rawBytes))
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	if err := newTpl.Execute(buf, params); err != nil {
		return nil, err
	}

	files["go/client.go"] = buf.Bytes()

	schemas := Models(spec)
	fmt.Printf("\n\nModels: %#v\n\n", schemas)

	for _, schema := range schemas {
		if schema.Scalar != nil {
			fmt.Printf("Scalar: %#v\n", schema.Scalar.Name)
		}

		if schema.Struct != nil {
			fmt.Printf("Struct: %#v\n", schema.Struct.Name)

			for _, p := range schema.Struct.Properties {
				fmt.Printf("Property: %s %s\n", p.Scalar.Name, p.Scalar.Type)
			}
		}
	}

	return files, nil
}

type operation struct {
	Name         string
	Description  string
	Method       string
	Path         string
	PathParams   []param
	QueryParams  []param
	HeaderParams []param
	CookieParams []param
}

func (o operation) OtherParams() []param {
	allParams := make([]param, 0)
	allParams = append(allParams, o.QueryParams...)
	allParams = append(allParams, o.HeaderParams...)
	allParams = append(allParams, o.CookieParams...)

	return allParams
}

func (o operation) HasPathParams() bool {
	return len(o.PathParams) > 0
}

func (o operation) HasQueryParams() bool {
	return len(o.QueryParams) > 0
}

func (o operation) HasHeaderParams() bool {
	return len(o.HeaderParams) > 0
}

func (o operation) HasCookieParams() bool {
	return len(o.CookieParams) > 0
}

func (o operation) HasOtherParams() bool {
	return o.HasQueryParams() || o.HasHeaderParams() || o.HasCookieParams()
}

func (o operation) JoinPathParams() []string {
	decls := make([]string, 0)

	for _, eachParam := range o.PathParams {
		if eachParam.Type == "path" {
			decls = append(decls, fmt.Sprintf("%s %s", eachParam.OriginalName, DetectType(eachParam.Schema)))
		}
	}

	return decls
}

func (o operation) ToArgs() string {
	args := o.JoinPathParams()

	if o.HasOtherParams() {
		args = append(args, fmt.Sprintf("params %sParams", o.Name))
	}

	return strings.Join(args, ", ")
}

type param struct {
	Name         string
	Type         string
	OriginalName string
	Template     string
	Schema       oas3.Schema
}

func DetectType(s oas3.Schema) string {
	if s.SchemaValue.Type == "string" {
		return "string"
	}

	if s.SchemaValue.Type == "integer" {
		if s.SchemaValue.Format == "int32" {
			return "int32"
		}

		if s.SchemaValue.Format == "int64" {
			return "int64"
		}

		return "int"
	}

	if s.SchemaValue.Type == "array" {
		return "[]interface{}"
	}

	return "interface{}"
}

func ConvertToString(variable string, s oas3.Schema) string {
	kind := DetectType(s)

	switch kind {
	case "int", "int32", "int64":
		return fmt.Sprintf("strconv.FormatInt(%s, 10)", variable)
	case "float", "float32", "float64":
		return fmt.Sprintf("strconv.FormatFloat(%s)", variable)
	case "string":
		return variable
	}

	return fmt.Sprintf("string(%s)", variable)
}

func Functions(o oas3.Spec) []Function {
	return nil
}

func SchemaToModel(propertyName string, propertyValue oas3.Schema) Types {
	newType := Types{}

	switch propertyValue.Type {
	case "integer":
		newScalar := Scalar{
			Name: propertyName,
			Tag:  propertyName,
		}

		if propertyValue.Format == "int32" {
			newScalar.Type = "int32"
		} else if propertyValue.Format == "int64" {
			newScalar.Type = "int64"
		} else {
			newScalar.Type = "int"
		}

		if stringutil.Contains(propertyValue.Required, propertyName) {
			newScalar.Type = "*" + newScalar.Type
		}

		newType.Scalar = &newScalar

	case "number":
		newScalar := Scalar{
			Name: propertyName,
			Tag:  propertyName,
		}

		if propertyValue.Format == "float" {
			newScalar.Type = "float32"
		} else if propertyValue.Format == "double" {
			newScalar.Type = "float64"
		} else {
			newScalar.Type = "float"
		}

		if stringutil.Contains(propertyValue.Required, propertyName) {
			newScalar.Type = "*" + newScalar.Type
		}

		newType.Scalar = &newScalar
	case "boolean":
		newScalar := Scalar{
			Name: propertyName,
			Type: "bool",
			Tag:  propertyName,
		}

		if stringutil.Contains(propertyValue.Required, propertyName) {
			newScalar.Type = "*" + newScalar.Type
		}

		newType.Scalar = &newScalar
	case "string":
		newScalar := Scalar{
			Name: propertyName,
			Type: "string",
			Tag:  propertyName,
		}

		if stringutil.Contains(propertyValue.Required, propertyName) {
			newScalar.Type = "*" + newScalar.Type
		}

		newType.Scalar = &newScalar
	case "object":
		props := make([]Types, 0)

		for k, v := range propertyValue.Properties {
			props = append(props, SchemaToModel(k, v))
		}

		newType.Struct = &Struct{
			Name:       propertyName,
			Properties: props,
			Tag:        propertyName,
		}
	}

	return newType
}

func Models(spec *oas3.Spec) []Types {
	types := make([]Types, 0)

	for schemaName, schemaValue := range spec.Components.Schemas {
		types = append(types, SchemaToModel(schemaName, schemaValue))
	}

	return types
}

type Types struct {
	Array  *Types
	Scalar *Scalar
	Struct *Struct
}

type Class Struct

type Module Package

type Package struct {
	Name string
}

type Struct struct {
	Implements string
	Extends    string
	Tag        string
	Name       string
	Properties []Types
	Functions  []Function
}

type Scalar struct {
	Name string
	Type string
	Tag  string
}

type Function struct {
	Name           string
	RequiredParams []param
	OptionalParams []param
}
