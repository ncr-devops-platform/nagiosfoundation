// +build windows

package nagiosfoundation

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func uint16SliceToString(array []uint16) string {
	var end int

	for end = 0; array[end] != 0; end++ {
	}

	return syscall.UTF16ToString(array[:end])
}

func showHelpOsConstrained() {
	fmt.Println("\nNote: Process names in Windows are not case sensitive.")
}

func isProcessRunningOsConstrained(name string) bool {
	retval := false

	handle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err == nil {
		defer windows.CloseHandle(handle)

		var entry windows.ProcessEntry32
		entry.Size = uint32(unsafe.Sizeof(entry))

		err = windows.Process32First(handle, &entry)

		for err == nil && retval == false {
			exeName := uint16SliceToString(entry.ExeFile[0:len(entry.ExeFile)])
			//fmt.Println("Entry:", exeName, "| Match:", name)
			retval = strings.EqualFold(name, exeName)

			err = windows.Process32Next(handle, &entry)
		}
	}

	return retval
}
