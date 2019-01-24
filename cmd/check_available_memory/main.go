// +build windows
package main

import (
	"github.com/jkerry/nagiosfoundation/lib/app/nagiosfoundation/check_available_memory"
)

func main() {
	nagiosfoundation.CheckAvailableMemory()
}
