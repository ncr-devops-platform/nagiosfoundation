package nagiosfoundation

import (
	"fmt"
	"strconv"
	
	"github.com/shirou/gopsutil/host"

)

func dealwithErr(err error) {
	if err != nil {
		fmt.Println(err)
		//os.Exit(-1)
	}
}

func getHostUptime() uint64 {
	hostStat, err := host.Info()
	dealwithErr(err)
	return hostStat.Uptime
}

// CheckUptime gathers information about the host uptime.
func CheckUptime(warning string) (string, int) {
	var msg string
	var retCode int
	var checkStateText, msgString string

	switch {
	default:
		uptime := getHostUptime()
		warn,_ := strconv.Atoi(warning)
		switch {
		case uptime >= uint64(warn):
			checkStateText = statusTextCritical
			msgString = fmt.Sprint("Warning: the system uptime is greater than the specified value. Uptime is: %s",strconv.FormatUint(uptime,10))
			retCode = 1
		case uptime < uint64(warn):
			checkStateText = statusTextOK
			msgString = fmt.Sprint("System uptime is below the specified value. Uptime is,: %s",strconv.FormatUint(uptime,10))
			retCode = 0
		}
	}

	msg, _ = resultMessage("CheckUptime", checkStateText, msgString)
	return msg, retCode
}
