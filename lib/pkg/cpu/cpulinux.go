// +build !windows

package cpu

func getCPULoadOsConstrained() (float64, error) {
	return getCPULoadLinux()
}
