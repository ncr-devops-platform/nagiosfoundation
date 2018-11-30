// +build windows

package nagiosfoundation

import (
	"flag"
	"fmt"
	"os"

	"github.com/shirou/gopsutil/winservices"
)

func CheckServiceRunning() {
	var serviceName = flag.String("service_name", "", "the name of the service to check")
	flag.Parse()
	s, _ := winservices.NewService(*serviceName)
	var ss winservices.ServiceStatus
	ss, _ = s.QueryStatus()
	var msg string
	var retcode int
	if ss.State != 4 {
		msg = fmt.Sprintf("CheckServiceRunning CRITICAL - %s not in a running state", *serviceName)
		retcode = 2
	} else {
		msg = fmt.Sprintf("CheckServiceRunning OK- %s in a running state", *serviceName)
		retcode = 0
	}
	fmt.Println(msg)
	os.Exit(retcode)
}
