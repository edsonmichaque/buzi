package golang

import (
	"errors"

	"github.com/edsonmichaque/umbeluzi/types"
)

type modelGenerator struct{}

func (m modelGenerator) Generate(spec *types.Spec) ([]byte, error) {
	objs := make(map[string]Type, 0)
	for name, schema := range spec.Components.Schemas {
		if schema.Type == "object" {
			nk := Type{
				Name:       name,
				properties: make(map[string]string, 0),
			}

			for k, v := range schema.Properties {
				if v.Type == "int" {
					nk.properties[k] = T{
						Kind: "int",
					}
				}

				if v.Type == "string" {
					nk.properties[k] = T{
						Kind: "string",
					}
				}

				if v.Type == "boolean" {
					nk.properties[k] = T{
						Kind: "bool",
					}
				}

				if v.Type == "number" {
					nk.properties[k] = T{
						Kind: "float64",
					}
				}
			}
		}
	}

	return nil, errors.New("buzi: not implemented")
}

var x = map[X]string{
	{Type: "integer"}: "int",
	{Type: "boolean"}: "bool",
	{Type: "string"}:  "string",
	{Type: "number"}:  "float64",

	{Type: "array", Items: Items{Type: "int"}}:     "[]int `json:\"\"`",
	{Type: "array", Items: Items{Type: "boolean"}}: "[]bool",
	{Type: "array", Items: Items{Type: "string"}}:  "[]string",
	{Type: "array", Items: Items{Type: "number"}}:  "[]number",
}

type X struct {
	Type  string
	Items Items
}

type Items struct {
	Type string
}

type Type struct {
	Name       string
	properties map[string]T
}

type T struct {
	Kind        string
	Annotations map[string]string
}
