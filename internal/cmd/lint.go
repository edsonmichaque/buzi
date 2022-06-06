package cmd

import (
	"github.com/spf13/cobra"
)

func Lint() *Command {
	newCobra := cobra.Command{
		Use: "boa",
		RunE: func(cmd *cobra.Command, args []string) error {
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
