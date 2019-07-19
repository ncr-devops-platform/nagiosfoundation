package nagiosfoundation

import (
	"fmt"
	
	"github.com/shirou/gopsutil/host"
	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/nagiosformatters"
)

//Just a tad bit of error handling

func dealwithErr(err error) {
	if err != nil {
		fmt.Println(err)
		//os.Exit(-1)
	}
}

//getHostUptime simply returns the host uptime as a uint64

func getHostUptime() uint64 {
	hostStat, err := host.Info()
	dealwithErr(err)
	return hostStat.Uptime
}

// CheckUptime gathers information about the host uptime.
func CheckUptime(checkType string, warning, critical int, metricName string) (string, int) {
	
	const checkName = "Checkuptime"
	
	var msg string
	var retcode int
	
	uptime := getHostUptime()
	
	msg, retcode = nagiosformatters.GreaterFormatNagiosCheck(checkName, float64(uptime), float64(warning), float64(critical), metricName)
	
	return msg, retcode
}