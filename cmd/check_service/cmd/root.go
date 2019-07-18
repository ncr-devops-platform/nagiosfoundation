package cmd

import (
	"fmt"
	"os"

	"github.com/ncr-devops-platform/nagiosfoundation/cmd/initcmd"
	"github.com/ncr-devops-platform/nagiosfoundation/lib/app/nagiosfoundation"
	"github.com/spf13/cobra"
)

const serviceManagerFlag = "manager"
const currentStateWantedFlag = "current_state"

var state, user, manager string
var currentStateWanted bool

// Execute runs the root command
func Execute() {
	var name string

	var rootCmd = &cobra.Command{
		Use:   "check_service",
		Short: "Determine the status of a service.",
		Long: `Perform various checks for a service. These checks depend on the options
given and the --name (-n) option is always required.` + getHelpOsConstrained(),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.ParseFlags(os.Args)

			msg, retcode := nagiosfoundation.CheckService(name, state, user, currentStateWanted, manager)

			fmt.Println(msg)
			os.Exit(retcode)
		},
	}

	initcmd.AddVersionCommand(rootCmd)

	const nameFlag = "name"
	rootCmd.Flags().StringVarP(&name, nameFlag, "n", "", "service name")
	rootCmd.MarkFlagRequired(nameFlag)

	addFlagsOsConstrained(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
