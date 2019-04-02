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
	var user, group string

	var rootCmd = &cobra.Command{
		Use:   "check_user_group",
		Short: "Determine if a user and/or group is on a system.",
		Long: `Checks for the existence of a user, a group, or if a user exists and
belongs to a group. At least one flag must be provided.

- The --user (-u) flag, checks for the existence of the user.
- The --group (-g) flag checks for the existence of a group.
- The --user (-u) and --group (-g) flag together check for the
  existence of the user, then checks that the user belongs in
  the group.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.ParseFlags(os.Args)

			if user == "" && group == "" {
				cmd.Help()
			} else {
				msg, retval := nagiosfoundation.CheckUserGroup(user, group)

				fmt.Println(msg)
				os.Exit(retval)
			}
		},
	}

	initcmd.AddVersionCommand(rootCmd)

	rootCmd.Flags().StringVarP(&user, "user", "u", "", "user name")
	rootCmd.Flags().StringVarP(&group, "group", "g", "", "group name")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
