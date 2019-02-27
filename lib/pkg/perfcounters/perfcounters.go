// +build windows

package perfcounters

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

type PerformanceCounter struct {
	Name  string
	Value float64
}

type PowerShell struct {
	powerShell string
}

func New() *PowerShell {
	ps, _ := exec.LookPath("powershell.exe")
	return &PowerShell{
		powerShell: ps,
	}
}

func (p *PowerShell) Execute(args ...string) (stdOut string, stdErr string, err error) {
	args = append([]string{"-NoProfile", "-NonInteractive"}, args...)
	cmd := exec.Command(p.powerShell, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	stdOut, stdErr = stdout.String(), stderr.String()
	return
}

// ReadPerformanceCounter reads a performance counter
func ReadPerformanceCounter(counter string, pollingAttempts int, pollingDelay int) (PerformanceCounter, error) {
	// in amd64 the pdh.dll usage isn't playing nice. We're going to use powershell directly and text parsing
	var perfcounter PerformanceCounter

	perfcounter.Name = counter
	perfcounter.Value = 0

	var command string
	command = fmt.Sprintf("Write-Output (Get-Counter -Counter \"%s\" -SampleInterval %d -MaxSamples %d |\n", counter, pollingDelay, pollingAttempts) +
		"Select-Object -ExpandProperty CounterSamples |\n" +
		"Select-Object -ExpandProperty CookedValue |\n" +
		"Measure-Object -Average).Average"

	if glog.V(2) {
		glog.Infof("Generated powershell performance monitor command:\n%s\n", command)
	}

	posh := New()
	stdout, _, err := posh.Execute(command)

	if glog.V(2) {
		glog.Infof("powershell output: \n\n %v", stdout)
	}

	if err != nil {
		if glog.V(1) {
			glog.Errorf("Error running powershell script to retrieve performance counter values: %v", err)
		}

		return perfcounter, err
	}

	trimmed_stdout := strings.TrimSpace(stdout)
	avgValue, err := strconv.ParseFloat(trimmed_stdout, 64)

	if err != nil {
		if glog.V(1) {
			glog.Errorf("Could not parse %s to float64: %v", stdout, err)
		}

		return perfcounter, err
	}

	perfcounter.Value = avgValue

	return perfcounter, nil

}
