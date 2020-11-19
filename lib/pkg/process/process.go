package process

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// GetProcessesByName returns FileInfo for stat files and pids for a given process name
func GetProcessesByName(name string) ([]GeneralInfo, error) {
	svc := processByNameHandlers{
		open: os.Open,
		close: func(f *os.File) error {
			return f.Close()
		},
		readDir: func(f *os.File, entries int) ([]os.FileInfo, error) {
			return f.Readdir(entries)
		},
		getPidName: getPidNameWithHandler,
		readFile:   ioutil.ReadFile,
	}

	return getProcessesByNameWithHandlers(svc, name)
}

func getProcessesByNameWithHandlers(svc processByNameHandlers, name string) ([]GeneralInfo, error) {
	var errorReturn error
	matchingEntries := make([]GeneralInfo, 0)

	dir, err := svc.open("/proc")
	if err != nil {
		matchingEntries = nil
		errorReturn = err
	}

	defer svc.close(dir)

	var procEntries []os.FileInfo
	if errorReturn == nil {
		procEntries, err = svc.readDir(dir, 0)

		if err != nil {
			matchingEntries = nil
			errorReturn = err
		}
	}

	if errorReturn == nil {
		for _, procEntry := range procEntries {
			// Skip entries that aren't directories
			if !procEntry.IsDir() {
				continue
			}

			// Skip entries that aren't numbers
			pid, err := strconv.Atoi(procEntry.Name())
			if err != nil {
				continue
			}

			if procName, _ := svc.getPidName(svc.readFile, pid); procName == name {
				matchingEntries = append(matchingEntries, GeneralInfo{
					PID:      pid,
					StatFile: procEntry,
				})
			}
		}
	}

	return matchingEntries, errorReturn
}

func getPidNameWithHandler(readFile func(string) ([]byte, error), pid int) (string, error) {
	procFile := fmt.Sprintf("/proc/%d/stat", pid)
	procDataBytes, err := readFile(procFile)
	if err != nil {
		return "", err
	}

	procData := string(procDataBytes)

	procNameStart := strings.IndexRune(procData, '(') + 1
	procNameEnd := strings.IndexRune(procData, ')')

	if procNameStart >= procNameEnd {
		return "", errors.New("Could not parse process name")
	}

	procName := procData[procNameStart:procNameEnd]

	return procName, nil
}

type processByNameHandlers struct {
	open       func(string) (*os.File, error)
	close      func(*os.File) error
	readDir    func(*os.File, int) ([]os.FileInfo, error)
	getPidName func(readFile func(string) ([]byte, error), pid int) (string, error)
	readFile   func(string) ([]byte, error)
}

// GeneralInfo represent process information, such as PID and stat file
type GeneralInfo struct {
	PID      int
	StatFile os.FileInfo
}
