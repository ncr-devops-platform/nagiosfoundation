package main

import (
	nf "github.com/ncr-devops-platform/nagiosfoundation/lib/app/nagiosfoundation"
	nagiosfoundation "github.com/ncr-devops-platform/nagiosfoundation/lib/app/nagiosfoundation/check_service"
)

func main() {
	nf.CheckExecutableVersion()
	nagiosfoundation.CheckService()
}
