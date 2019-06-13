package cmd

import (
	"fmt"
	"os"

	"github.com/ncr-devops-platform/nagiosfoundation/cmd/initcmd"
	"github.com/spf13/cobra"
)

// Execute runs the root command
func Execute(apiCheckHTTP func(string, bool, int, string, string, string) (string, int)) int {
	var exitCode int
	var url string
	var redirect bool
	var timeout int
	var format string
	var path string
	var expectedValue string
	var rootCmd = &cobra.Command{
		Use:   "check_http",
		Short: "Check the response code of an http request.",
		Long:  `Perform an HTTP get request and assert whether it is OK, warning or critical.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.ParseFlags(os.Args)
			msg, retval := apiCheckHTTP(url, redirect, timeout, format, path, expectedValue)

			fmt.Println(msg)
			exitCode = retval
		},
	}

	initcmd.AddVersionCommand(rootCmd)

	rootCmd.Flags().StringVarP(&url, "url", "u", "http://127.0.0.1", "the URL to check")
	rootCmd.Flags().BoolVarP(&redirect, "redirect", "r", false, "follow redirects?")
	rootCmd.Flags().IntVarP(&timeout, "timeout", "t", 15, "timeout in seconds")
	rootCmd.Flags().StringVarP(&format, "format", "f", "", "The expected response format: json")
	rootCmd.Flags().StringVarP(&path, "path", "p", "", "The path in the return value data to test against the expected value")
	rootCmd.Flags().StringVarP(&expectedValue, "expectedValue", "e", "", "The expected response data value")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		exitCode = 1
	}

	return exitCode
}
