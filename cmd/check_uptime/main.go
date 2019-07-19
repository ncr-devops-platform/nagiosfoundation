package main

import (
	"os"

	"github.com/ncr-devops-platform/nagiosfoundation/cmd/check_uptime/cmd"
	"github.com/ncr-devops-platform/nagiosfoundation/lib/app/nagiosfoundation"
)

func main() {
	os.Exit(cmd.Execute(nagiosfoundation.CheckUptime))
}
