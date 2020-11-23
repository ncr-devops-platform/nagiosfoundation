// +build !windows

package memory

import (
	"errors"
	"strings"
	"testing"

	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/process"
)

func Test_handlePSOutput(t *testing.T) {
	type args struct {
		output string
		err    error
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "EmptyOutput",
			args: args{
				output: "",
			},
			want:    0.0,
			wantErr: true,
		},
		{
			name: "ShellError",
			args: args{
				output: "abc\n 2.23343",
				err:    errors.New("TEST"),
			},
			want:    0.0,
			wantErr: true,
		},
		{
			name: "InvalidFormat",
			args: args{
				output: "abc\nqwe",
			},
			want:    0.0,
			wantErr: true,
		},
		{
			name: "ValidCase",
			args: args{
				output: "abc\n 1.234567",
			},
			want:    1.234567,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handlePSOutput(tt.args.output, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("handlePSOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handlePSOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getProcessMemoryPercentageOsContrained(t *testing.T) {
	tests := []struct {
		name                       string
		getProcessesByNameOverride func(string) ([]process.GeneralInfo, error)
		execBashOverride           func(string) ([]byte, error)
		want                       float64
		wantErr                    bool
	}{
		{
			name: "GetProcessesError",
			getProcessesByNameOverride: func(_ string) ([]process.GeneralInfo, error) {
				return nil, errors.New("TEST")
			},
			want:    0.0,
			wantErr: true,
		},
		{
			name: "NoProcesses",
			getProcessesByNameOverride: func(_ string) ([]process.GeneralInfo, error) {
				return []process.GeneralInfo{}, nil
			},
			want:    0.0,
			wantErr: false,
		},
		{
			name: "CommandError",
			getProcessesByNameOverride: func(_ string) ([]process.GeneralInfo, error) {
				return []process.GeneralInfo{
					process.GeneralInfo{
						PID: 123,
					},
				}, nil
			},
			execBashOverride: func(cmd string) ([]byte, error) {
				return nil, nil
			},
			want:    0.0,
			wantErr: true,
		},
		{
			name: "SingleProcess",
			getProcessesByNameOverride: func(_ string) ([]process.GeneralInfo, error) {
				return []process.GeneralInfo{
					process.GeneralInfo{
						PID: 123,
					},
				}, nil
			},
			execBashOverride: func(cmd string) ([]byte, error) {
				return []byte("%mem\n 1.23456"), nil
			},
			want:    1.23456,
			wantErr: false,
		},
		{
			name: "MultipleProcesses",
			getProcessesByNameOverride: func(_ string) ([]process.GeneralInfo, error) {
				return []process.GeneralInfo{
					process.GeneralInfo{
						PID: 123,
					},
					process.GeneralInfo{
						PID: 456,
					},
				}, nil
			},
			execBashOverride: func(cmd string) ([]byte, error) {
				if strings.Contains(cmd, "123") {
					return []byte("%mem\n 1.25"), nil
				}
				return []byte("%mem\n 1.5"), nil
			},
			want:    2.75,
			wantErr: false,
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

			oldExecBash := execBash
			if tt.execBashOverride != nil {
				execBash = tt.execBashOverride
				defer func() {
					execBash = oldExecBash
				}()
			}

			got, err := getProcessMemoryPercentageOsContrained("test")
			if (err != nil) != tt.wantErr {
				t.Errorf("getProcessMemoryPercentageOsContrained() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getProcessMemoryPercentageOsContrained() = %v, want %v", got, tt.want)
			}
		})
	}
}
