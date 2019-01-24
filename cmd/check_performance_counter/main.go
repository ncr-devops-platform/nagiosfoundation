// +build windows
package main

import (
	"github.com/jkerry/nagiosfoundation/lib/app/nagiosfoundation/check_performance_counter"
)

func main() {
	nagiosfoundation.CheckPerformanceCounter()
}
