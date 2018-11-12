package perfcounters

import (
	"errors"
	"fmt"
	"time"
	"unsafe"

	"github.com/lxn/win"
)

type PerformanceCounter struct {
	Name  string
	Value float64
}

// ReadPerformanceCounter reads a performance counter
func ReadPerformanceCounter(counter string, pollingAttempts int, pollingDelay int) (PerformanceCounter, error) {

	var queryHandle win.PDH_HQUERY
	var counterHandle win.PDH_HCOUNTER
	var perfcounter PerformanceCounter

	ret := win.PdhOpenQuery(0, 0, &queryHandle)
	if ret != win.ERROR_SUCCESS {
		return perfcounter, errors.New("Unable to open query through DLL call")
	}

	// test path
	ret = win.PdhValidatePath(counter)
	if ret == win.PDH_CSTATUS_BAD_COUNTERNAME {
		return perfcounter, errors.New("Unable to fetch counter (this is unexpected)")
	}

	ret = win.PdhAddEnglishCounter(queryHandle, counter, 0, &counterHandle)
	if ret != win.ERROR_SUCCESS {
		return perfcounter, fmt.Errorf("unable to add process counter. Error code is %x", ret)
	}

	ret = win.PdhCollectQueryData(queryHandle)
	if ret != win.ERROR_SUCCESS {
		return perfcounter, fmt.Errorf("got an error: 0x%x", ret)
	}

	var collect = func(samples int, waitTime int) float64 {
		var data []float64
		for index := 0; index < samples; index++ {
			ret = win.PdhCollectQueryData(queryHandle)
			if ret == win.ERROR_SUCCESS {

				var bufSize uint32
				var bufCount uint32
				var size uint32 = uint32(unsafe.Sizeof(win.PDH_FMT_COUNTERVALUE_ITEM_DOUBLE{}))
				var emptyBuf [1]win.PDH_FMT_COUNTERVALUE_ITEM_DOUBLE // need at least 1 addressable null ptr.

				ret = win.PdhGetFormattedCounterArrayDouble(counterHandle, &bufSize, &bufCount, &emptyBuf[0])
				if ret == win.PDH_MORE_DATA {
					filledBuf := make([]win.PDH_FMT_COUNTERVALUE_ITEM_DOUBLE, bufCount*size)
					ret = win.PdhGetFormattedCounterArrayDouble(counterHandle, &bufSize, &bufCount, &filledBuf[0])
					if ret == win.ERROR_SUCCESS {
						for i := 0; i < int(bufCount); i++ {
							c := filledBuf[i]
							data = append(data, c.FmtValue.DoubleValue)
						}
					}
				}
			}

			time.Sleep(time.Duration(waitTime) * time.Second)
		}
		//first value is always 0. Not sure why, but popping it
		var finalSample []float64
		var x float64
		x, finalSample = data[0], data[1:]
		if x != 0 {
			finalSample = data
		}
		var total float64 = 0
		for _, value := range finalSample {
			total += value
		}

		return (total / float64(len(finalSample)))
	}

	perfcounter.Name = counter
	perfcounter.Value = collect(3, 1)

	return perfcounter, nil

}
