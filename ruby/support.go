package ruby

import "github.com/edsonmichaque/buzi"

type support struct {
	makefileTemplate     string
	gosumTemplate        string
	gomodTemplate        string
	readmeTemplate       string
	contributingTemplate string
}

func (m support) Apply(_ map[string]string, manifest *buzi.Manifest) ([]buzi.File, error) {
	files := []buzi.File{
		{
			Path: "go.mod",
		},
		{
			Path: "go.sum",
		},
		{
			Path: "LICENSE",
		},
		{
			Path: "README.md",
		},
		{
			Path: "CONTRIBUTING.md",
		},
	}

	return files, nil
}
