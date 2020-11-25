// +build !windows

package cpu

import (
	"errors"
	"testing"

	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/process"
)

func Test_parseTopSamples(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "EmptyLines",
			args: args{
				lines: []string{},
			},
			want:    0.0,
			wantErr: false,
		},
		{
			name: "InvalidElemCount",
			args: args{
				lines: []string{"test"},
			},
			want:    0.0,
			wantErr: true,
		},
		{
			name: "NonFloatCPUValue",
			args: args{
				lines: []string{"1 2 3 4 5 6 7 8 CPU 10 11 12"},
			},
			want:    0.0,
			wantErr: true,
		},
		{
			name: "ValidSingleLine",
			args: args{
				lines: []string{"1 2 3 4 5 6 7 8 55.5 10 11 12"},
			},
			want:    55.5,
			wantErr: false,
		},
		{
			name: "ValidMultiLine",
			args: args{
				lines: []string{
					"1 2 3 4 5 6 7 8 55.5 10 11 12",
					"1 2 3 4 5 6 7 8 150.0 10 11 12",
				},
			},
			want:    205.5,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTopSamples(tt.args.lines)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTopSamples() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseTopSamples() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleTopOutput(t *testing.T) {
	type args struct {
		output string
		err    error
	}
	tests := []struct {
		name                string
		args                args
		getCPUCountOverride func() int
		want                float64
		wantErr             bool
	}{
		{
			name: "ShellError",
			args: args{
				output: "",
				err:    errors.New("TEST"),
			},
			getCPUCountOverride: func() int {
				return 1
			},
			want:    0.0,
			wantErr: true,
		},
		{
			name: "NoPIDInOutput",
			args: args{
				output: "abc\nqwe\nrty\n",
			},
			getCPUCountOverride: func() int {
				return 1
			},
			want:    0.0,
			wantErr: false,
		},
		{
			name: "OneSample",
			args: args{
				output: " PID\n 1 2 3 4 5 6 7 8 99.9 10 11 12\n",
			},
			getCPUCountOverride: func() int {
				return 1
			},
			want:    0.0,
			wantErr: false,
		},
		{
			name: "InvalidSampleLine",
			args: args{
				output: " PID\nPID\n 1 2 3 4 5 6 7 8 CPU 10 11 12\n",
			},
			getCPUCountOverride: func() int {
				return 1
			},
			want:    0.0,
			wantErr: true,
		},
		{
			name: "ValidSingleSampleLine",
			args: args{
				output: " PID\n 1 2 3 4 5 6 7 8 99.9 10 11 12\nPID\n 1 2 3 4 5 6 7 8 99.9 10 11 12\n",
			},
			getCPUCountOverride: func() int {
				return 1
			},
			want:    99.9,
			wantErr: false,
		},
		{
			name: "ValidMultipleSampleLines",
			args: args{
				output: " PID\n 1 2 3 4 5 6 7 8 99.9 10 11 12\nPID\n 1 2 3 4 5 6 7 8 50 10 11 12\n 1 2 3 4 5 6 7 8 150 10 11 12\n\nPID\n 1 2 3 4 5 6 7 8 20 10 11 12\n 1 2 3 4 5 6 7 8 30 10 11 12\n",
			},
			getCPUCountOverride: func() int {
				return 2
			},
			want:    62.5,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		oldGetCPUCount := getCPUCount
		if tt.getCPUCountOverride != nil {
			getCPUCount = tt.getCPUCountOverride
			defer func() {
				getCPUCount = oldGetCPUCount
			}()
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := handleTopOutput(tt.args.output, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleTopOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handleTopOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getProcessCPULoad(t *testing.T) {
	type args struct {
		processInfo []process.GeneralInfo
	}
	tests := []struct {
		name                string
		args                args
		execBashOverride    func(string) ([]byte, error)
		getCPUCountOverride func() int
		want                float64
		wantErr             bool
	}{
		{
			name: "NoProcesses",
			args: args{
				processInfo: []process.GeneralInfo{},
			},
			want:    0.0,
			wantErr: false,
		},
		{
			name: "ShellError",
			args: args{
				processInfo: []process.GeneralInfo{
					process.GeneralInfo{},
				},
			},
			execBashOverride: func(_ string) ([]byte, error) {
				return nil, errors.New("TEST")
			},
			want:    0.0,
			wantErr: true,
		},
		{
			name: "ValidCase",
			args: args{
				processInfo: []process.GeneralInfo{
					process.GeneralInfo{
						PID: 1,
					},
				},
			},
			execBashOverride: func(_ string) ([]byte, error) {
				return []byte("PID\n1 2 3\nPID\n1 2 3 4 5 6 7 8 50.0 10 11 12"), nil
			},
			getCPUCountOverride: func() int {
				return 1
			},
			want:    50.0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldExecBash := execBash
			if tt.execBashOverride != nil {
				execBash = tt.execBashOverride
				defer func() {
					execBash = oldExecBash
				}()
			}

			oldGetCPUCount := getCPUCount
			if tt.getCPUCountOverride != nil {
				getCPUCount = tt.getCPUCountOverride
				defer func() {
					getCPUCount = oldGetCPUCount
				}()
			}

			got, err := getProcessCPULoad(tt.args.processInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("getProcessCPULoad() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getProcessCPULoad() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getProcessCPULoadOsConstrained(t *testing.T) {
	type args struct {
		processName        string
		perCoreCalculation bool
	}
	tests := []struct {
		name                       string
		args                       args
		getProcessesByNameOverride func(string) ([]process.GeneralInfo, error)
		want                       float64
		wantErr                    bool
	}{
		{
			name: "NoProcessesWithCore",
			args: args{
				perCoreCalculation: true,
			},
			getProcessesByNameOverride: func(_ string) ([]process.GeneralInfo, error) {
				return []process.GeneralInfo{}, nil
			},
			want:    0.0,
			wantErr: false,
		},
		{
			name: "NoProcessesWithoutCore",
			args: args{
				perCoreCalculation: false,
			},
			getProcessesByNameOverride: func(_ string) ([]process.GeneralInfo, error) {
				return []process.GeneralInfo{}, nil
			},
			want:    0.0,
			wantErr: false,
		},
		{
			name: "GetProcessesError",
			args: args{
				perCoreCalculation: true,
			},
			getProcessesByNameOverride: func(_ string) ([]process.GeneralInfo, error) {
				return nil, errors.New("TEST")
			},
			want:    0.0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldGetProcessesByName := getProcessesByName
			if tt.getProcessesByNameOverride != nil {
				getProcessesByName = tt.getProcessesByNameOverride
				defer func() {
					getProcessesByName = oldGetProcessesByName
				}()
			}

			got, err := getProcessCPULoadOsConstrained(tt.args.processName, tt.args.perCoreCalculation)
			if (err != nil) != tt.wantErr {
				t.Errorf("getProcessCPULoadOsConstrained() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getProcessCPULoadOsConstrained() = %v, want %v", got, tt.want)
			}
		})
	}
}
