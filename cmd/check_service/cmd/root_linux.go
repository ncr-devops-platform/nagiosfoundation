// +build !windows

package cmd

import "github.com/spf13/cobra"

func getHelpOsConstrained() string {
	return `

For Linux, the only check done is for a running state. Both the --name (-n) and
-manager (-m) options must be specified and the service is only checked
to see if it is running.
`
}

func addFlagsOsConstrained(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&currentStateWanted, currentStateWantedFlag, "c", false, "output the service state in nagios output")
	cmd.Flags().StringVarP(&manager, serviceManagerFlag, "m", "systemd", "the name of local service manager. Allowed options are: systemd")
}
