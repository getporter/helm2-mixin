package main

import (
	"github.com/deislabs/porter-helm/pkg/helm"
	"github.com/spf13/cobra"
)

func buildStatusCommand(m *helm.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Execute the status functionality of this mixin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Status()
		},
	}
	return cmd
}
