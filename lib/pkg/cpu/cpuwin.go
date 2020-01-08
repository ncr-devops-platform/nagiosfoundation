// +build windows

package cpu

import (
	"github.com/StackExchange/wmi"
)

type Win32_PerfFormattedData_PerfOS_Processor struct {
	PercentProcessorTime				uint64
}

func getCPULoadOsConstrained() (float64, error) {
	var dst []Win32_PerfFormattedData_PerfOS_Processor
	q := wmi.CreateQuery(&dst, "WHERE name='_Total'")
	err := wmi.Query(q, &dst)
	if(err != nil){
		return 0.0, err
	}
	return (float64)(dst[0].PercentProcessorTime), nil
}
