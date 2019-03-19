// +build !windows

package memory

// GetFreeMemoryOsConstrained returns the amount of available memory.
func getFreeMemoryOsConstrained() uint64 {
	return getMemInfoEntryFromFile("/proc/meminfo", "MemAvailable")
}
