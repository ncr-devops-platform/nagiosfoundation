package cmd

import (
	"fmt"
	"os"

	"github.com/ncr-devops-platform/nagiosfoundation/cmd/initcmd"
	"github.com/spf13/cobra"
)

// Execute runs the root command
func Execute(apiCheckUptime func(string) (string, int)) int {
	var exitCode int
	var warning string

	var rootCmd = &cobra.Command{
		Use:   "check_uptime",
		Short: "Check the uptime of a system.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.ParseFlags(os.Args)
			msg, retval := apiCheckUptime(warning)

			fmt.Println(msg)
			exitCode = retval
		},
	}

	initcmd.AddVersionCommand(rootCmd)

	rootCmd.Flags().StringVarP(&warning, "pattern", "w", "259200", "The desired warning level in seconds for uptime, defaults to 72 hours (259200s)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		exitCode = 1
	}

	return exitCode
}
