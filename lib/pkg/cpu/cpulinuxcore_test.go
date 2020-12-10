// +build !windows

package cpu

import (
	"errors"
	"reflect"
	"testing"

	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/process"
)

func Test_parsePIDStatLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    *pidStatLine
		wantErr bool
	}{
		{
			name: "InvalidElemCount",
			args: args{
				line: " 1 2 3 4 5 6 7 8 9 ",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "NonIntTID",
			args: args{
				line: " 1 2 3 TID 5 6 7 8 9 10",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "NonIntTID2",
			args: args{
				line: " 1 2 3 4.2 5 6 7 8 9 10",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "NonFloatCPUValue",
			args: args{
				line: " 1 2 3 4 5 6 7 CPU 9 10",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "NonIntCPUCore",
			args: args{
				line: " 1 2 3 4 5 6 7 8 9.2 10",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "NonIntCPUCore2",
			args: args{
				line: " 1 2 3 4 5 6 7 8 CORE 10",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Valid",
			args: args{
				line: " 1 2 3 177 5 6 7 8.5 3 10",
			},
			want: &pidStatLine{
				TID:  177,
				CPU:  8.5,
				Core: 3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePIDStatLine(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePIDStatLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePIDStatLine() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parsePIDStatSample(t *testing.T) {
	type args struct {
		sampleLines []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*pidStatLine
		wantErr bool
	}{
		{
			name: "NoLines",
			args: args{
				sampleLines: []string{},
			},
			want:    []*pidStatLine{},
			wantErr: false,
		},
		{
			name: "ParseError",
			args: args{
				sampleLines: []string{
					"1 2 3 a 5 6 7 b c 10",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "SingleLine",
			args: args{
				sampleLines: []string{
					"1 2 3 4 5 6 7 8 9 10",
				},
			},
			want: []*pidStatLine{
				&pidStatLine{
					TID:  4,
					CPU:  8.0,
					Core: 9,
				},
			},
			wantErr: false,
		},
		{
			name: "MultipleLines",
			args: args{
				sampleLines: []string{
					"1 2 3 4 5 6 7 8 9 10",
					"10 20 30 40 50 60 70 80 90 100",
				},
			},
			want: []*pidStatLine{
				&pidStatLine{
					TID:  4,
					CPU:  8.0,
					Core: 9,
				},
				&pidStatLine{
					TID:  40,
					CPU:  80.0,
					Core: 90,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePIDStatSample(tt.args.sampleLines)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePIDStatSample() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePIDStatSample() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSampleRanges(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		name string
		args args
		want map[int]*sampleStartEnd
	}{
		{
			name: "InvalidSample",
			args: args{
				lines: []string{
					"Linux 3.33.3-3333.3.3.e17.x33_33",
					"#  Time UID",
					"#  Time UID",
				},
			},
			want: map[int]*sampleStartEnd{},
		},
		{
			name: "BrokenValidSample",
			args: args{
				lines: []string{
					"Linux 3.33.3-3333.3.3.e17.x33_33",
					"#  Time UID",
					"#  Time UID",
					"#  Time UID",
					"1606200000 0 117 0 16.0 80.0 0.00 96.00 1 command",
					"1606200000 0 0 117 8.0 40.0 0.00 48.00 1 command",
					"1606200000 0 0 118 5.0 36.0 0.00 41.00 0 command",
					"1606200000 0 0 119 3.0 4.0 0.00 7.00 3 command",
				},
			},
			want: map[int]*sampleStartEnd{
				2: &sampleStartEnd{
					Start: 4,
					End:   8,
				},
			},
		},
		{
			name: "BrokenValidSample2",
			args: args{
				lines: []string{
					"Linux 3.33.3-3333.3.3.e17.x33_33",
					"#  Time UID",
					"1606200000 0 117 0 16.0 80.0 0.00 96.00 1 command",
					"1606200000 0 0 117 8.0 40.0 0.00 48.00 1 command",
					"1606200000 0 0 118 5.0 36.0 0.00 41.00 0 command",
					"1606200000 0 0 119 3.0 4.0 0.00 7.00 3 command",
					"#  Time UID",
				},
			},
			want: map[int]*sampleStartEnd{
				0: &sampleStartEnd{
					Start: 2,
					End:   6,
				},
			},
		},
		{
			name: "SingleSample",
			args: args{
				lines: []string{
					"Linux 3.33.3-3333.3.3.e17.x33_33",
					"#  Time UID",
					"1606200000 0 117 0 16.0 80.0 0.00 96.00 1 command",
					"1606200000 0 0 117 8.0 40.0 0.00 48.00 1 command",
					"1606200000 0 0 118 5.0 36.0 0.00 41.00 0 command",
					"1606200000 0 0 119 3.0 4.0 0.00 7.00 3 command",
				},
			},
			want: map[int]*sampleStartEnd{
				0: &sampleStartEnd{
					Start: 2,
					End:   6,
				},
			},
		},
		{
			name: "MultipleSamples",
			args: args{
				lines: []string{
					"Linux 3.33.3-3333.3.3.e17.x33_33",
					"#  Time UID",
					"1606200000 0 117 0 16.0 80.0 0.00 96.00 1 command",
					"1606200000 0 0 117 8.0 40.0 0.00 48.00 1 command",
					"1606200000 0 0 118 5.0 36.0 0.00 41.00 0 command",
					"1606200000 0 0 119 3.0 4.0 0.00 7.00 3 command",
					"#  Time UID",
					"1606200001 0 117 0 8.0 40.0 0.00 48.00 1 command",
					"1606200001 0 0 117 4.0 20.0 0.00 24.00 1 command",
					"1606200001 0 0 118 2.5 18.0 0.00 20.50 0 command",
					"1606200001 0 0 119 1.5 2.0 0.00 3.50 3 command",
				},
			},
			want: map[int]*sampleStartEnd{
				0: &sampleStartEnd{
					Start: 2,
					End:   6,
				},
				1: &sampleStartEnd{
					Start: 7,
					End:   11,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSampleRanges(tt.args.lines); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSampleRanges() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNonEmptyLines(t *testing.T) {
	type args struct {
		input []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "MultipleEmptyLines",
			args: args{
				input: []string{
					"",
					"\n",
					"\t   ",
					"abc",
				},
			},
			want: []string{
				"abc",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNonWhiteSpaceLines(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNonEmptyLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateCPUCoreUsage(t *testing.T) {
	type args struct {
		pidStatLines []*pidStatLine
	}
	tests := []struct {
		name string
		args args
		want map[int]float64
	}{
		{
			name: "NoLines",
			args: args{
				pidStatLines: []*pidStatLine{},
			},
			want: map[int]float64{},
		},
		{
			name: "SingleCore",
			args: args{
				pidStatLines: []*pidStatLine{
					&pidStatLine{
						CPU:  25.0,
						Core: 0,
						TID:  0,
					},
					&pidStatLine{
						CPU:  15.0,
						Core: 0,
						TID:  177,
					},
					&pidStatLine{
						CPU:  5.0,
						Core: 0,
						TID:  178,
					},
				},
			},
			want: map[int]float64{
				0: 20.0,
			},
		},
		{
			name: "MultipleCores",
			args: args{
				pidStatLines: []*pidStatLine{
					&pidStatLine{
						CPU:  39.0,
						Core: 0,
						TID:  0,
					},
					&pidStatLine{
						CPU:  15.0,
						Core: 0,
						TID:  177,
					},
					&pidStatLine{
						CPU:  5.0,
						Core: 1,
						TID:  178,
					},
					&pidStatLine{
						CPU:  15.5,
						Core: 1,
						TID:  179,
					},
					&pidStatLine{
						CPU:  3.5,
						Core: 1,
						TID:  180,
					},
				},
			},
			want: map[int]float64{
				0: 15.0,
				1: 24.0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateCPUCoreUsage(tt.args.pidStatLines); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calculateCPUCoreUsage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateCPUCoreAverage(t *testing.T) {
	type args struct {
		input map[int]map[int]float64
	}
	tests := []struct {
		name string
		args args
		want map[int]float64
	}{
		{
			name: "NoData",
			args: args{
				input: map[int]map[int]float64{},
			},
			want: map[int]float64{},
		},
		{
			name: "ValidCase",
			args: args{
				input: map[int]map[int]float64{
					1: map[int]float64{
						0: 25.5,
						1: 33.3,
					},
					2: map[int]float64{
						1: 25.5,
						2: 33.3,
					},
					3: map[int]float64{
						0: 33.5,
						3: 25.5,
					},
				},
			},
			want: map[int]float64{
				0: 29.5,
				1: 29.4,
				2: 33.3,
				3: 25.5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateCPUCoreAverage(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calculateCPUCoreAverage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCPUCoreUsage(t *testing.T) {
	type args struct {
		output   string
		shellErr error
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "ShellError",
			args: args{
				shellErr: errors.New("TEST"),
			},
			want:    0.0,
			wantErr: true,
		},
		{
			name: "ParseError",
			args: args{
				output: "Linux 3.33.3-3333.3.3.e17.x33_33\n#  Time UID\n#  Time UID\n#  Time UID\n1606200000 0 0 TID 16.0 80.0 0.00 96.00 1 command",
			},
			want:    0.0,
			wantErr: true,
		},
		{
			name: "ValidCase",
			args: args{
				output: "Linux 3.33.3-3333.3.3.e17.x33_33\n\n#  Time UID\n1606200000 0 117 0 16.0 80.0 0.00 96.00 1 command\n1606200000 0 0 117 8.0 40.0 0.00 48.00 1 command\n1606200000 0 0 118 5.0 36.0 0.00 41.00 0 command\n1606200000 0 0 119 3.0 4.0 0.00 7.00 3 command\n",
			},
			want:    48.0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getCPUCoreUsage(tt.args.output, tt.args.shellErr)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCPUCoreUsage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getCPUCoreUsage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getProcessCoreCPULoad(t *testing.T) {
	type args struct {
		processInfo []process.GeneralInfo
	}
	tests := []struct {
		name             string
		args             args
		execBashOverride func(string) ([]byte, error)
		want             float64
		wantErr          bool
	}{
		{
			name: "ShellError",
			args: args{
				processInfo: []process.GeneralInfo{
					process.GeneralInfo{
						PID: 1,
					},
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
				return []byte("Linux 3.33.3-3333.3.3.e17.x33_33\n\n#  Time UID\n1606200000 0 117 0 16.0 80.0 0.00 96.00 1 command\n1606200000 0 0 117 8.0 40.0 0.00 48.00 1 command\n1606200000 0 0 118 5.0 36.0 0.00 41.00 0 command\n1606200000 0 0 119 3.0 4.0 0.00 7.00 3 command\n"), nil
			},
			want:    48.0,
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

			got, err := getProcessCoreCPULoad(tt.args.processInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("getProcessCoreCPULoad() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getProcessCoreCPULoad() = %v, want %v", got, tt.want)
			}
		})
	}
}
