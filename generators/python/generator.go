package python

import (
	"github.com/edsonmichaque/buzi"
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
