package memory

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func availableMemoryErrorText(description string, expected, actual uint64) string {
	return fmt.Sprintf("%s. Expected Result: %d Actual Result %d",
		description, expected, actual)
}

func TestAvailableMemory(t *testing.T) {
	totalMemoryReturned := uint64(10000)
	freeMemoryReturned := uint64(6000)
	getTotalMemory := func() uint64 { return totalMemoryReturned }
	getTotalMemoryZero := func() uint64 { return uint64(0) }
	getFreeMemory := func() uint64 { return freeMemoryReturned }
	getFreeMemoryZero := func() uint64 { return uint64(0) }

	actualResult := getTotalMemoryWithHandler(getTotalMemory)
	if actualResult != totalMemoryReturned {
		t.Error(availableMemoryErrorText("getTotalMemoryWithHandler returned memory not correct",
			totalMemoryReturned, actualResult))
	}

	actualResult = getFreeMemoryWithHandler(getFreeMemory)
	if actualResult != freeMemoryReturned {
		t.Error(availableMemoryErrorText("getFreeMemoryWithHandler returned memory not correct",
			freeMemoryReturned, actualResult))
	}

	actualResult = getUsedMemoryWithHandlers(getTotalMemory, getFreeMemory)
	expectedResult := totalMemoryReturned - freeMemoryReturned
	if actualResult != expectedResult {
		t.Error(availableMemoryErrorText("getUsedMemoryWithHandlers did not calculate used memory properly",
			expectedResult, actualResult))
	}

	actualResult = getUsedMemoryWithHandlers(getTotalMemoryZero, getFreeMemory)
	expectedResult = 0
	if actualResult != expectedResult {
		t.Error(availableMemoryErrorText("getUsedMemoryWithHandlers did not properly handle a total memory of 0",
			expectedResult, actualResult))
	}

	actualResult = getUsedMemoryWithHandlers(getTotalMemory, getFreeMemoryZero)
	expectedResult = 0
	if actualResult != expectedResult {
		t.Error(availableMemoryErrorText("getUsedMemoryWithHandlers did not properly handle free memory of 0",
			expectedResult, actualResult))
	}

	actualResult = getUsedMemoryPercentageWithHandlers(getTotalMemory, getFreeMemory)
	expectedResult = uint64(float64((totalMemoryReturned - freeMemoryReturned)) / float64(totalMemoryReturned) * 100)
	if actualResult != expectedResult {
		t.Error(availableMemoryErrorText("getUsedMemoryPercentageWithHandlers did not calculated used memory percentage properly",
			expectedResult, actualResult))
	}

	actualResult = getUsedMemoryPercentageWithHandlers(getTotalMemoryZero, getFreeMemory)
	expectedResult = 0
	if actualResult != expectedResult {
		t.Error(availableMemoryErrorText("getUsedMemoryPercentageWithHandlers did not properly handle a total memory of 0",
			expectedResult, actualResult))
	}

	actualResult = getUsedMemoryPercentageWithHandlers(getTotalMemory, getFreeMemoryZero)
	expectedResult = 0
	if actualResult != expectedResult {
		t.Error(availableMemoryErrorText("getUsedMemoryPercentageWithHandlers did not properly handle free memory of 0",
			expectedResult, actualResult))
	}

	// Non-DI calls should at least return some numbers
	actualResult = GetTotalMemory()
	if actualResult == 0 {
		t.Error("GetTotalMemory did return > 0")
	}

	actualResult = GetFreeMemory()
	if actualResult == 0 {
		t.Error("GetFreeMemory did return > 0")
	}

	actualResult = GetUsedMemory()
	if actualResult == 0 {
		t.Error("GetUsedMemory did return > 0")
	}

	actualResult = GetUsedMemoryPercentage()
	if actualResult == 0 {
		t.Error("GetUsedMemoryPercentage did return > 0")
	}
}

func TestLinuxGetAvailableMemory(t *testing.T) {
	type parameters struct {
		filename string
		match    string
	}

	memInfoContents := `This string
has a valid
MemAvailable:    123456 kB
line`

	availableMemory := getMemInfoEntryFromReader(strings.NewReader(memInfoContents), "MemAvailable")
	if availableMemory != 123456*1024 {
		t.Error("Failed to parse memory amount from valid line")
	}

	memInfoContents = `This string has
no valid entry
to parse
`
	availableMemory = getMemInfoEntryFromReader(strings.NewReader(memInfoContents), "MemAvailable")
	if availableMemory != 0 {
		t.Error("Invalid memory info string failed to generate error")
	}

	memInfoContents = `This string has
MemAvailable:    breakme kB
to parse
	`
	availableMemory = getMemInfoEntryFromReader(strings.NewReader(memInfoContents), "MemAvailable")
	if availableMemory != 0 {
		t.Error("Invalid memory info string failed to generate error")
	}

	memInfoContents = `This string has
MemAvailable:   28446744073709551615 kB
to parse
	`
	availableMemory = getMemInfoEntryFromReader(strings.NewReader(memInfoContents), "MemAvailable")
	if availableMemory != 0 {
		t.Error("Invalid memory info string failed to generate error")
	}

	availableMemory = getMemInfoEntryFromFile(os.Args[0], "")
	if availableMemory != 0 {
		t.Error("Invalid memory info file and match string failed to generate error")
	}
}
