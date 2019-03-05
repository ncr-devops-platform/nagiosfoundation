// +build windows

package nagiosfoundation

import (
	"fmt"
	"os"

	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/perfcounters"
)

// CheckPerformanceCounter executes CheckPerformanceCounterWitHandler(),
// passing it the OS constranted ReadPerformanceCounter() function, prints
// the returned message and exits with the returned exit code.
func CheckPerformanceCounter() {
	msg, retval := CheckPerformanceCounterWithHandler(perfcounters.ReadPerformanceCounter)

	fmt.Println(msg)
	os.Exit(retval)
}
