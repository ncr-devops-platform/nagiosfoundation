package nagiosfoundation

import (
	"os"
	"testing"
	"time"
)

func TestCheckUptime(t *testing.T) {
	forever := time.Duration(24000 * time.Hour)

	// Deep down in the github.com/shirou/gopsutil/host package is
	// is code to override the uptime retrieval location using the
	// HOST_PROC environment variable. Set this to cause it to fail.
	// Test for failure first because on success, the package caches
	// the boot time used for uptime calculation.
	os.Setenv("HOST_PROC", "nosuchfile")
	msg, retcode := CheckUptime("", forever, forever, "uptime_name")
	if len(msg) == 0 {
		t.Error("Message not populated for uptime CRITICAL")
	}

	if retcode != 2 {
		t.Error("Incorrect return code for uptime CRITICAL")
	}

	os.Unsetenv("HOST_PROC")
	msg, retcode = CheckUptime("", forever, forever, "uptime_name")
	if len(msg) == 0 {
		t.Error("Message not populated for uptime OK")
	}

	if retcode != 0 {
		t.Error("Incorrect return code for uptime OK")
	}
}
