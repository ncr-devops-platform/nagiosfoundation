package nagiosfoundation

import (
	"fmt"
	"time"
		
	"github.com/shirou/gopsutil/host"
	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/nagiosformatters"
)

// CheckUptime gathers information about the host uptime.
func CheckUptime(checkType string, warning, critical time.Duration, metricName string) (string, int) {
	
	const checkName = "CheckUptime"
	
	var msg string
	var retcode int
	
	uptime, err := host.Uptime()
	
	if err != nil {
		msg, _ = resultMessage(checkName, statusTextCritical, fmt.Sprintf("Failed to determine uptime %s", err.Error()))
		retcode = 2
	} else {
		msg, retcode = nagiosformatters.GreaterFormatNagiosCheck(checkName, float64(uptime), float64(warning.Seconds()), float64(critical.Seconds()), metricName)
	}
	
	return msg, retcode
}