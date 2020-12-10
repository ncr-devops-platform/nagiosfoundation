// +build windows

package cpu

import (
	"errors"
	"reflect"
	"testing"
)

func Test_getProcessCPULoadOsConstrained(t *testing.T) {
	type args struct {
		processName string
	}
	tests := []struct {
		name                      string
		args                      args
		getWin32ProcessesOverride func(string, interface{}) error
		getCPUCountOverride       func() int
		want                      float64
		wantErr                   bool
	}{
		{
			name: "Win32Error",
			args: args{
				processName: "test",
			},
			getWin32ProcessesOverride: func(_ string, _ interface{}) error {
				return errors.New("TEST")
			},
			getCPUCountOverride: func() int {
				return 1
			},
			want:    0.0,
			wantErr: true,
		},
		{
			name: "NoProcesses",
			args: args{
				processName: "test",
			},
			getWin32ProcessesOverride: func(_ string, _ interface{}) error {
				return nil
			},
			getCPUCountOverride: func() int {
				return 1
			},
			want:    0.0,
			wantErr: false,
		},
		{
			name: "Match",
			args: args{
				processName: "test",
			},
			getWin32ProcessesOverride: func(_ string, dst interface{}) error {
				ptrValue := reflect.ValueOf(dst).Elem()
				process := win32_PerfFormattedData_PerfProc_Process{
					Name:                 "test",
					PercentProcessorTime: 15,
				}
				process2 := win32_PerfFormattedData_PerfProc_Process{
					Name:                 "testio",
					PercentProcessorTime: 15,
				}
				process3 := win32_PerfFormattedData_PerfProc_Process{
					Name:                 "test#17",
					PercentProcessorTime: 30,
				}
				ptrValue.Set(reflect.Append(ptrValue, reflect.ValueOf(process)))
				ptrValue.Set(reflect.Append(ptrValue, reflect.ValueOf(process2)))
				ptrValue.Set(reflect.Append(ptrValue, reflect.ValueOf(process3)))
				return nil
			},
			getCPUCountOverride: func() int {
				return 1
			},
			want:    45.0,
			wantErr: false,
		},
		{
			name: "CpuCount",
			args: args{
				processName: "test",
			},
			getWin32ProcessesOverride: func(_ string, dst interface{}) error {
				ptrValue := reflect.ValueOf(dst).Elem()
				process := win32_PerfFormattedData_PerfProc_Process{
					Name:                 "test",
					PercentProcessorTime: 20,
				}
				ptrValue.Set(reflect.Append(ptrValue, reflect.ValueOf(process)))
				return nil
			},
			getCPUCountOverride: func() int {
				return 4
			},
			want:    5.0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		oldGetWin32Processes := getWin32Processes
		if tt.getWin32ProcessesOverride != nil {
			getWin32Processes = tt.getWin32ProcessesOverride
			defer func() {
				getWin32Processes = oldGetWin32Processes
			}()
		}

		oldGetCPUCount := getCPUCount
		if tt.getCPUCountOverride != nil {
			getCPUCount = tt.getCPUCountOverride
			defer func() {
				getCPUCount = oldGetCPUCount
			}()
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := getProcessCPULoadOsConstrained(tt.args.processName, false)
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
