package nagiosfoundation

import (
	"errors"
	"testing"
)

func TestCheckProcessCPUWithHandler(t *testing.T) {
	type args struct {
		warning               int
		critical              int
		processName           string
		metricName            string
		perCoreCalculation    bool
		processCPUCoreHandler func(string, bool) (float64, error)
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
				processCPUCoreHandler: func(_ string, _ bool) (float64, error) {
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
				processCPUCoreHandler: func(_ string, _ bool) (float64, error) {
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
				processCPUCoreHandler: func(_ string, _ bool) (float64, error) {
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
				processCPUCoreHandler: func(_ string, _ bool) (float64, error) {
					return 50, nil
				},
			},
			want: 0,
		},
		{
			name: "NoCPUHandler",
			args: args{
				warning:               50,
				critical:              99,
				processName:           "test",
				metricName:            "test",
				processCPUCoreHandler: nil,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, retcode := CheckProcessCPUWithHandler(tt.args.warning, tt.args.critical, tt.args.processName, tt.args.metricName, tt.args.perCoreCalculation, tt.args.processCPUCoreHandler)
			if retcode != tt.want {
				t.Errorf("CheckProcessCPUWithHandler() got = %v, want %v", retcode, tt.want)
			}
		})
	}
}
