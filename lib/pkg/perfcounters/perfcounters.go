package perfcounters

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// PerformanceCounter is a struct that contains the name
// of the performance counter and the value retrieved
// from the performance counter.
type PerformanceCounter struct {
	Name  string
	Value float64
}

type powerShellService interface {
	Execute(...string) (string, string, error)
}

type powerShell struct {
	powerShell string

	command func(string, ...string) *exec.Cmd
}

func newPowerShell(lookPath func(string) (string, error), command func(string, ...string) *exec.Cmd) *powerShell {
	ps, _ := lookPath("powershell.exe")

	return &powerShell{
		powerShell: ps,
		command:    command,
	}
}

func (p *powerShell) Execute(args ...string) (stdOut string, stdErr string, err error) {
	args = append([]string{"-NoProfile", "-NonInteractive"}, args...)

	cmd := p.command(p.powerShell, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	stdOut, stdErr = stdout.String(), stderr.String()
	return
}

// ReadPerformanceCounterWithHandler reads a performance counter
func ReadPerformanceCounterWithHandler(poshService powerShellService, counter string, pollingAttempts int, pollingDelay int) (PerformanceCounter, error) {
	// in amd64 the pdh.dll usage isn't playing nice. We're going to use powershell directly and text parsing
	var perfcounter PerformanceCounter

	perfcounter.Name = counter
	perfcounter.Value = 0

	if poshService == nil {
		return PerformanceCounter{}, errors.New("No Powershell Execute service")
	}

	var command string
	command = fmt.Sprintf("Write-Output (Get-Counter -Counter \"%s\" -SampleInterval %d -MaxSamples %d |\n", counter, pollingDelay, pollingAttempts) +
		"Select-Object -ExpandProperty CounterSamples |\n" +
		"Select-Object -ExpandProperty CookedValue |\n" +
		"Measure-Object -Average).Average"

	stdout, _, err := poshService.Execute(command)

	if err != nil {
		return perfcounter, err
	}

	trimmedStdout := strings.TrimSpace(stdout)

	if trimmedStdout == "" {
		return perfcounter, fmt.Errorf("No data returned for %s so this counter probably doesn't exist", counter)
	}

	avgValue, err := strconv.ParseFloat(trimmedStdout, 64)

	if err != nil {
		return perfcounter, fmt.Errorf("Error processing for counter (%s): %s", counter, err)
	}

	perfcounter.Value = avgValue

	return perfcounter, nil
}

// ReadPerformanceCounter reads a performance counter
func ReadPerformanceCounter(counter string, pollingAttempts int, pollingDelay int) (PerformanceCounter, error) {
	return ReadPerformanceCounterWithHandler(newPowerShell(exec.LookPath, exec.Command), counter, pollingAttempts, pollingDelay)
}
