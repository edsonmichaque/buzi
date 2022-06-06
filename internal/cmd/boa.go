package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/edsonmichaque/buzi/codegen"
	jsonParser "github.com/edsonmichaque/buzi/parsers/json"
	yamlParser "github.com/edsonmichaque/buzi/parsers/yaml"
	"github.com/edsonmichaque/buzi/providers/python"
	"github.com/edsonmichaque/buzi/readers"
	"github.com/edsonmichaque/buzi/readers/http"
	"github.com/edsonmichaque/buzi/readers/local"
	"github.com/edsonmichaque/go-openapi/oas3"
	"github.com/spf13/cobra"
)

type CmdError struct {
	ExitCode int
	err      error
}

func (c CmdError) Error() string {
	return c.err.Error()
}

type Command struct {
	cobra *cobra.Command
}

func (c Command) Execute() error {
	return c.cobra.Execute()
}

func New() *Command {
	newCobra := cobra.Command{
		Use: "boa",
		RunE: func(cmd *cobra.Command, args []string) error {

			specPath, err := cmd.Flags().GetString("from-file")
			if err != nil {
				return CmdError{
					ExitCode: 1,
					err:      err,
				}
			}

			var reader readers.Reader
			if strings.HasPrefix(specPath, "http://") || strings.HasPrefix(specPath, "https://") {
				reader, err = http.New(specPath)
				if err != nil {
					return CmdError{
						ExitCode: 1,
						err:      err,
					}
				}
			} else if strings.HasPrefix(specPath, "file://") {
				reader, err = local.New(strings.TrimPrefix(specPath, "file://"))
				if err != nil {
					return CmdError{
						ExitCode: 1,
						err:      err,
					}
				}
			} else {
				reader, err = local.New(specPath)
				if err != nil {
					return CmdError{
						ExitCode: 1,
						err:      err,
					}
				}
			}

			newYamlParser := yamlParser.New(reader)

			spec, err := oas3.New(newYamlParser)
			if err != nil {
				newJsonParser := jsonParser.New(reader)

				spec, err = oas3.New(newJsonParser)
				if err != nil {
					return CmdError{
						ExitCode: 1,
						err:      err,
					}
				}
			}

			for path, pathItem := range spec.Paths {
				ctx := oas3.ValidationContext{Path: path}
				if err := pathItem.Validate(ctx); err != nil {
					fmt.Println(err.Error())
				}
			}

			p := python.New()

			defs, err := codegen.NewDef(spec)
			if err != nil {
				return CmdError{
					ExitCode: 1,
					err:      err,
				}
			}

			options := codegen.OptionSet(
				map[string]interface{}{
					"gem-name":     "thisgem",
					"package-name": "thisgem",
					"module-name":  "ThisModule",
				},
			)

			codegen.Generate(p, *defs, options)

			files, err := codegen.Generate(p, *defs, options)
			if err != nil {
				return CmdError{
					ExitCode: 1,
					err:      err,
				}
			}

			for _, file := range files {

				newPath := path.Join("tmp", "ruby", file.Path)

				fmt.Println(newPath)

				if err := os.MkdirAll(path.Dir(newPath), 0755); err != nil {
					return CmdError{
						ExitCode: 1,
						err:      err,
					}
				}

				f, err := os.OpenFile(newPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					return CmdError{
						ExitCode: 1,
						err:      err,
					}
				}

				if _, err := f.Write(file.Content); err != nil {
					return CmdError{
						ExitCode: 1,
						err:      err,
					}
				}

				f.Close()
			}

			return nil
		},
	}

	newCobra.Flags().StringP("from-file", "f", "openapi.yml", "spec")
	newCobra.Flags().StringP("http-method", "M", "get", "HTTP method to use")
	newCobra.Flags().StringArrayP("http-header", "H", nil, "HTTP headers to send")
	newCobra.Flags().StringArrayP("query-param", "Q", nil, "HTTP query parameters")

	newCmd := Command{cobra: &newCobra}

	return &newCmd
}
