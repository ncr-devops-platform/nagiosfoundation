package cpu

import (
	"errors"
	"testing"
)

func TestCpu(t *testing.T) {
	var cpuData [2]string
	var testError error

	goodData := [2]string{
		"cpu  10 10 10 00 10 10 10 10 10 10",
		"cpu  20 20 20 100 30 20 20 20 20 20",
	}

	badBeforeData := [2]string{
		"cpu  235445 1470 abc 15137277 201523 0 3938 0 0 0",
		"cpu  235445 1470 74483 15137277 201523 0 3938 0 0 0",
	}

	badAfterData := [2]string{
		"cpu  235445 1470 74481 15137277 201523 0 3938 0 0 0",
		"cpu  235445 1470 abc 15137277 201523 0 3938 0 0 0",
	}

	emptyData := [2]string{"", ""}

	shortData := [2]string{"cpu 10", "cpu 20"}

	var counter int

	testCPUData := func() (string, error) {
		data := cpuData[counter]
		counter++

		return data, testError
	}

	cpuData = goodData
	testError = nil
	counter = 0
	usage, err := getCPULoadLinuxWithHandler(testCPUData)

	if err != nil {
		t.Error("Good CPU data should not return an error")
	}

	if usage != 50 {
		t.Error("Good CPU data should give usage of 50")
	}

	cpuData = goodData
	testError = errors.New("Test Error")
	counter = 0
	usage, err = getCPULoadLinuxWithHandler(testCPUData)

	if err == nil {
		t.Error("Error getting CPU data should return an error")
	}

	if usage != 0 {
		t.Error("Error getting CPU data should give usage of 0")
	}

	cpuData = badAfterData
	testError = nil
	counter = 0
	usage, err = getCPULoadLinuxWithHandler(testCPUData)

	if err == nil {
		t.Error("Getting error on second CPU check should return an error")
	}

	if usage != 0 {
		t.Error("Bad CPU data should give usage of 0")
	}

	cpuData = badBeforeData
	testError = nil
	counter = 0
	usage, err = getCPULoadLinuxWithHandler(testCPUData)

	if err == nil {
		t.Error("Unparseable CPU data should return an error")
	}

	if usage != 0 {
		t.Error("Bad CPU data should give usage of 0")
	}

	cpuData = emptyData
	testError = nil
	counter = 0
	usage, err = getCPULoadLinuxWithHandler(testCPUData)

	if err == nil {
		t.Error("Empty CPU data should return an error")
	}

	if usage != 0 {
		t.Error("Empty CPU data should give usage of 0")
	}

	cpuData = shortData
	testError = nil
	counter = 0
	usage, err = getCPULoadLinuxWithHandler(testCPUData)

	if err == nil {
		t.Error("Short line of CPU data should return an error")
	}

	if usage != 0 {
		t.Error("Short line of CPU data should give usage of 0")
	}

	if _, err := getStats(nil); err == nil {
		t.Error("No stats data handler should return an error")
	}

	if _, err := getCPULoadLinuxWithHandler(nil); err == nil {
		t.Error("No stats data handler should return an error")
	}
	// Execute to at least make sure there's no panic
	getStatsDataService()
}

func Test_average(t *testing.T) {
	type args struct {
		input []float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Average",
			args: args{
				input: []float64{
					5,
					10,
					14.8,
					30,
				},
			},
			want: 14.95,
		},
		{
			name: "Zero",
			args: args{
				input: []float64{},
			},
			want: 0.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := average(tt.args.input); got != tt.want {
				t.Errorf("average() = %v, want %v", got, tt.want)
			}
		})
	}
}
