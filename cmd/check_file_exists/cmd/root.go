package cmd

import (
	"fmt"
	"os"

	"github.com/ncr-devops-platform/nagiosfoundation/cmd/initcmd"
	"github.com/ncr-devops-platform/nagiosfoundation/lib/app/nagiosfoundation"
	"github.com/spf13/cobra"
)

// Execute runs the root command
func Execute() {
	var pattern string
	var negate bool

	var rootCmd = &cobra.Command{
		Use:   "check_file_exists",
		Short: "Check for the existence of one or more files matching specific filepath or globbing patterns.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.ParseFlags(os.Args)
			msg, retval := nagiosfoundation.CheckFileExists(pattern, negate)

			fmt.Println(msg)
			os.Exit(retval)
		},
	}

	initcmd.AddVersionCommand(rootCmd)

	rootCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "Filepath or globbing pattern to check for one or more existing files")
	rootCmd.Flags().BoolVarP(&negate, "negate", "n", false, "If set, asserts filepath or globbing pattern should not match any existing file")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
