package cmd

import (
	"fmt"
	"os"

	"github.com/ncr-devops-platform/nagiosfoundation/lib/app/nagiosfoundation"
	"github.com/spf13/cobra"
)

// Execute runs the root command
func Execute() {
	var name, checkType string

	var rootCmd = &cobra.Command{
		Use:   "check_process",
		Short: "Determine if a process is running.",
		Long:  nagiosfoundation.GetHelpProcess(),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.ParseFlags(os.Args)
			nagiosfoundation.CheckProcess(name, checkType)
		},
	}

	nagiosfoundation.AddVersionCommand(rootCmd)

	const nameFlag = "name"
	rootCmd.Flags().StringVarP(&name, nameFlag, "n", "", "process name")
	rootCmd.MarkFlagRequired(nameFlag)
	rootCmd.Flags().StringVarP(&checkType, "type", "t", "running", "Supported types are \"running\" and \"notrunning\"")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
