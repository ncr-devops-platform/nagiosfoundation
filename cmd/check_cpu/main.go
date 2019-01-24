// +build windows
package main

import (
	"github.com/jkerry/nagiosfoundation/lib/app/nagiosfoundation/check_cpu"
)

func main() {
	nagiosfoundation.CheckCPU()
}
