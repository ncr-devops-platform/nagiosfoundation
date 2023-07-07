// +build windows

package cmd

import "github.com/spf13/cobra"

func getHelpOsConstrained() string {
	return `
	
Some examples:
  check_service.exe --name audiosrv
    Checks for the service to exist and shows the service state and user.
  check_service.exe --name audiosrv --state running
    Checks for the service in the running state.
  check_service.exe --name audiosrv --state running --user "NT AUTHORITY\LocalService"
    Checks for the service in the running state and running as user.
  check_service.exe --name audiosrv --user "NT AUTHORITY\LocalService"
    Checks for the service to exist and would be run as user.
`
}

func addFlagsOsConstrained(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&state, "state", "s", "", "the desired state of the service")
	cmd.Flags().StringVarP(&user, "user", "u", "", "the user the service should run as")
	cmd.Flags().BoolVarP(&currentStateWanted, currentStateWantedFlag, "c", false, "output the Windows service state in nagios output")
	cmd.Flags().BoolVarP(&useSvcMgr, serviceManagerFlag, "w", false, "Decides which Windows service manager to use. Default is false / \"Windows Management Instrumentation (wmi)\". Set to true to use \"Windows Service Manager\"")
}
