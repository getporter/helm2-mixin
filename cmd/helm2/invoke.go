package main

import (
	"get.porter.sh/mixin/helm2/pkg/helm2"
	"github.com/spf13/cobra"
)

func buildInvokeCommand(mixin *helm2.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invoke",
		Short: "Execute the invoke functionality of this mixin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return mixin.Execute()
		},
	}

	// Define a flag for --action so that its presence doesn't cause errors, but ignore it since it will
	// be derived from the unmarshaled payload sent to it
	var action string
	cmd.Flags().StringVar(&action, "action", "", "Custom action name to invoke.")

	return cmd
}
