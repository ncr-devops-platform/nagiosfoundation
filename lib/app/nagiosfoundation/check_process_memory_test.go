package nagiosfoundation

import (
	"errors"
	"testing"
)

func TestCheckProcessMemoryWithHandler(t *testing.T) {
	type args struct {
		warning       int
		critical      int
		processName   string
		metricName    string
		memoryHandler func(string) (float64, error)
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "WarningResponse",
			args: args{
				warning:     50,
				critical:    99,
				processName: "test",
				metricName:  "test",
				memoryHandler: func(_ string) (float64, error) {
					return 55, nil
				},
			},
			want: 1,
		},
		{
			name: "CriticalResponse",
			args: args{
				warning:     50,
				critical:    99,
				processName: "test",
				metricName:  "test",
				memoryHandler: func(_ string) (float64, error) {
					return 100, nil
				},
			},
			want: 2,
		},
		{
			name: "CriticalResponse2",
			args: args{
				warning:     50,
				critical:    99,
				processName: "test",
				metricName:  "test",
				memoryHandler: func(_ string) (float64, error) {
					return 0.0, errors.New("TEST")
				},
			},
			want: 2,
		},
		{
			name: "OkResponse",
			args: args{
				warning:     50,
				critical:    99,
				processName: "test",
				metricName:  "test",
				memoryHandler: func(_ string) (float64, error) {
					return 50, nil
				},
			},
			want: 0,
		},
		{
			name: "NoMemoryHandler",
			args: args{
				warning:       50,
				critical:      99,
				processName:   "test",
				metricName:    "test",
				memoryHandler: nil,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, retcode := CheckProcessMemoryWithHandler(tt.args.warning, tt.args.critical, tt.args.processName, tt.args.metricName, tt.args.memoryHandler)
			if retcode != tt.want {
				t.Errorf("CheckProcessMemoryWithHandler() got = %v, want %v", retcode, tt.want)
			}
		})
	}
}
