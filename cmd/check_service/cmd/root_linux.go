// +build !windows

package cmd

import "github.com/spf13/cobra"

func addFlagsOsConstrained(cmd *cobra.Command) {
	const managerFlag = "manager"

	cmd.Flags().StringVarP(&manager, managerFlag, "m", "", "the name of local service manager. Allowed options are: systemd")
	cmd.MarkFlagRequired(managerFlag)
}
