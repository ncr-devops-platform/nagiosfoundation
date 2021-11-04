package nagiosfoundation

import (
	"testing"
	"time"
)

func TestCheckUptime(t *testing.T) {
	forever := time.Duration(24000 * time.Hour)
	msg, retcode := CheckUptime("", forever, forever, "uptime_name")
	if len(msg) == 0 {
		t.Error("Message not populated for uptime OK")
	}

	if retcode != 0 {
		t.Error("Incorrect return code for uptime OK")
	}
}
