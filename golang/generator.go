package golang

import (
	"github.com/edsonmichaque/buzi"
)

const (
	ParamModule  = "module"
	ParamPackage = "package"
)

func Pipeline() []buzi.Generator {
	return []buzi.Generator{
		models{
			path: "models",
		},
		operations{},
		support{
			makefileTemplate: "Makefile.tpl",
			readmeTemplate:   "README.tpl",
		},
	}
}
