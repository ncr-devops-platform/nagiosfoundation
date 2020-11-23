// +build windows

package memory

import (
	"errors"
	"reflect"
	"testing"
)

func generateProcessInfoWithName(id, parentID uint32, memory uint64, name string) *processInfo {
	return &processInfo{
		Process: win32_process{
			Name:            name,
			ProcessID:       id,
			ParentProcessID: parentID,
			WorkingSetSize:  memory,
		},
		UsedInCalculation: false,
	}
}

func generateProcessInfo(id, parentID uint32, memory uint64) *processInfo {
	return generateProcessInfoWithName(id, parentID, memory, "test")
}
func Test_getMemoryUsedByPIDAndItsChildren(t *testing.T) {
	type args struct {
		data map[uint32]*processInfo
		pid  uint32
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "NoProcesses",
			args: args{
				data: map[uint32]*processInfo{},
				pid:  177,
			},
			want: 0,
		},
		{
			name: "OneProcess",
			args: args{
				data: map[uint32]*processInfo{
					1: generateProcessInfo(1, 1000, 5000),
				},
				pid: 1,
			},
			want: 5000,
		},
		{
			name: "SingleNestedChild",
			args: args{
				data: map[uint32]*processInfo{
					1: generateProcessInfo(1, 1000, 5000),
					2: generateProcessInfo(2, 1, 2000),
				},
				pid: 1,
			},
			want: 7000,
		},
		{
			name: "TwoNestedChildren",
			args: args{
				data: map[uint32]*processInfo{
					1: generateProcessInfo(1, 1000, 5000),
					2: generateProcessInfo(2, 1, 2000),
					3: generateProcessInfo(3, 2, 1000),
				},
				pid: 1,
			},
			want: 8000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMemoryUsedByPIDAndItsChildren(tt.args.data, tt.args.pid); got != tt.want {
				t.Errorf("getMemoryUsedByPIDAndItsChildren() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getMemoryUsedByProcessNameAndItsChildren(t *testing.T) {
	type args struct {
		data map[uint32]*processInfo
		name string
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "NoProcesses",
			args: args{
				data: map[uint32]*processInfo{
					1: generateProcessInfoWithName(1, 1000, 5000, "other-name"),
				},
				name: "test",
			},
			want: 0,
		},
		{
			name: "OneProcess",
			args: args{
				data: map[uint32]*processInfo{
					1: generateProcessInfoWithName(1, 1000, 5000, "test"),
				},
				name: "test",
			},
			want: 5000,
		},
		{
			name: "MultipleProcesses",
			args: args{
				data: map[uint32]*processInfo{
					1: generateProcessInfoWithName(1, 1000, 5000, "test"),
					2: generateProcessInfoWithName(1, 1000, 5000, "test"),
				},
				name: "test",
			},
			want: 10000,
		},
		{
			name: "SingleNestedChild",
			args: args{
				data: map[uint32]*processInfo{
					1: generateProcessInfoWithName(1, 1000, 5000, "test"),
					2: generateProcessInfoWithName(2, 1, 2000, "sub-test"),
				},
				name: "test",
			},
			want: 7000,
		},
		{
			name: "MultipleProcessesSingleNestedChild",
			args: args{
				data: map[uint32]*processInfo{
					1:  generateProcessInfoWithName(1, 1000, 5000, "test"),
					2:  generateProcessInfoWithName(2, 1, 2000, "sub-test"),
					10: generateProcessInfoWithName(10, 1000, 10000, "test"),
					20: generateProcessInfoWithName(20, 10, 20000, "sub-test-2"),
				},
				name: "test",
			},
			want: 37000,
		},
		{
			name: "TwoNestedChildren",
			args: args{
				data: map[uint32]*processInfo{
					1: generateProcessInfoWithName(1, 1000, 5000, "test"),
					2: generateProcessInfoWithName(2, 1, 2000, "sub-test"),
					3: generateProcessInfoWithName(3, 2, 1000, "sub-sub-test"),
				},
				name: "test",
			},
			want: 8000,
		},
		{
			name: "MultipleProcessesTwoNestedChildren",
			args: args{
				data: map[uint32]*processInfo{
					1:  generateProcessInfoWithName(1, 1000, 5000, "test"),
					2:  generateProcessInfoWithName(2, 1, 2000, "sub-test"),
					3:  generateProcessInfoWithName(3, 2, 1000, "sub-sub-test"),
					10: generateProcessInfoWithName(10, 1000, 15000, "test"),
					20: generateProcessInfoWithName(20, 10, 2600, "sub-test-2"),
					30: generateProcessInfoWithName(30, 20, 1400, "sub-sub-test-2"),
				},
				name: "test",
			},
			want: 27000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMemoryUsedByProcessNameAndItsChildren(tt.args.data, tt.args.name); got != tt.want {
				t.Errorf("getMemoryUsedByProcessNameAndItsChildren() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getProcessMemoryPercentageOsContrained(t *testing.T) {
	type args struct {
		processName string
	}
	tests := []struct {
		name                      string
		args                      args
		getWin32ProcessesOverride func(string, interface{}, ...interface{}) error
		want                      float64
		wantErr                   bool
	}{
		{
			name: "Win32Error",
			args: args{
				processName: "test",
			},
			getWin32ProcessesOverride: func(_ string, dst interface{}, _ ...interface{}) error {
				return errors.New("TEST")
			},
			want:    0.0,
			wantErr: true,
		},
		{
			name: "SingleProcess",
			args: args{
				processName: "test",
			},
			getWin32ProcessesOverride: func(_ string, dst interface{}, _ ...interface{}) error {
				ptrValue := reflect.ValueOf(dst).Elem()
				process := win32_process{
					Name:            "test.exe",
					ProcessID:       1,
					ParentProcessID: 2,
					WorkingSetSize:  5000,
				}
				ptrValue.Set(reflect.Append(ptrValue, reflect.ValueOf(process)))
				return nil
			},
			want:    5.0,
			wantErr: false,
		},
		{
			name: "MultipleProcesses",
			args: args{
				processName: "test",
			},
			getWin32ProcessesOverride: func(_ string, dst interface{}, _ ...interface{}) error {
				ptrValue := reflect.ValueOf(dst).Elem()
				process := win32_process{
					Name:            "test.exe",
					ProcessID:       1,
					ParentProcessID: 2,
					WorkingSetSize:  5000,
				}
				process2 := win32_process{
					Name:            "test.exe",
					ProcessID:       10,
					ParentProcessID: 0,
					WorkingSetSize:  7500,
				}
				ptrValue.Set(reflect.Append(ptrValue, reflect.ValueOf(process)))
				ptrValue.Set(reflect.Append(ptrValue, reflect.ValueOf(process2)))
				return nil
			},
			want:    12.5,
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
		oldGetTotalMemory := getTotalMemory
		getTotalMemory = func() uint64 {
			return 100000
		}
		defer func() {
			getTotalMemory = oldGetTotalMemory
		}()

		t.Run(tt.name, func(t *testing.T) {
			got, err := getProcessMemoryPercentageOsContrained(tt.args.processName)
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
