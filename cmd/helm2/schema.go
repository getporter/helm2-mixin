package main

import (
	"get.porter.sh/mixin/helm2/pkg/helm2"
	"github.com/spf13/cobra"
)

func buildSchemaCommand(m *helm2.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Print the json schema for the mixin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.PrintSchema()
		},
	}
	return cmd
}
