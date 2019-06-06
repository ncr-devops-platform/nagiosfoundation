package cmd

import (
	"fmt"
	"os"

	"github.com/ncr-devops-platform/nagiosfoundation/cmd/initcmd"
	"github.com/spf13/cobra"
)

// Execute runs the root command
func Execute(apiCheckFileExists func(string, bool) (string, int)) int {
	var exitCode int
	var pattern string
	var negate bool

	var rootCmd = &cobra.Command{
		Use:   "check_file_exists",
		Short: "Check for the existence of one or more files matching specific filepath or globbing patterns.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.ParseFlags(os.Args)
			msg, retval := apiCheckFileExists(pattern, negate)

			fmt.Println(msg)
			exitCode = retval
		},
	}

	initcmd.AddVersionCommand(rootCmd)

	rootCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "Filepath or globbing pattern to check for one or more existing files")
	rootCmd.Flags().BoolVarP(&negate, "negate", "n", false, "Asserts filepath or globbing pattern should NOT match any existing file")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		exitCode = 1
	}

	return exitCode
}
