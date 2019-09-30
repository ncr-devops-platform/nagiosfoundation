package nagiosfoundation

import (
	"testing"
)

func TestCheckMemory(t *testing.T) {
	testReturnValid := func() uint64 { return uint64(50) }
	testReturnZero := func() uint64 { return uint64(0) }

	// No "get memory" service passed
	msg, retcode := CheckMemoryWithHandler("", 85, 95, "available_memory_percent", nil)

	if retcode != 2 || msg == "" {
		t.Error("CheckMemoryWithHandler() failed to handle nil service")
	}

	// Valid memory service with flag defaults
	msg, retcode = CheckMemoryWithHandler("", 85, 95, "available_memory_percent", testReturnValid)

	if retcode != 0 || msg == "" {
		t.Error("CheckMemoryWithHandler() failed with valid GetFreeMemory() call")
	}

	// Valid memory service but service returns error
	msg, retcode = CheckMemoryWithHandler("", 85, 95, "available_memory_percent", testReturnZero)

	if retcode != 2 || msg == "" {
		t.Error("CheckMemoryWithHandler() failed with valid GetFreeMemory() call")
	}

	msg, retcode = CheckMemoryWithHandler("", 40, 95, "available_memory_percent", testReturnValid)

	if retcode != 1 || msg == "" {
		t.Error("CheckMemoryWithHandler() should have emitted WARNING")
	}

	msg, retcode = CheckMemoryWithHandler("", 85, 45, "available_memory_percent", testReturnValid)

	if retcode != 2 || msg == "" {
		t.Error("CheckMemoryWithHandler() should have emitted CRITICAL")
	}
}
