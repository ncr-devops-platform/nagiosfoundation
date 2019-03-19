// +build windows

package cmd

import "github.com/spf13/cobra"

func addFlagsOsConstrained(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&state, "state", "s", "", "the desired state of the service")
	cmd.Flags().StringVarP(&user, "user", "u", "", "the user the service should run as")
}
