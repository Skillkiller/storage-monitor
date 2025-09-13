package metrics

import (
	"reflect"
	"testing"
)

func TestParseDfOutput(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []FileSystemInfo
		wantErr bool
	}{
		{
			name: "Standard output",
			input: `
Filesystem      1K-blocks      Used Available Use% Mounted on
u408683-sub4   1073741824 825296768 248445056  77% /home
tmpfs            1024000    1000   1023001     1% /run
`,
			want: []FileSystemInfo{
				{
					Filesystem: "u408683-sub4",
					Blocks:     1073741824,
					Used:       825296768,
					Available:  248445056,
					UsePercent: 77,
					MountedOn:  "/home",
				},
				{
					Filesystem: "tmpfs",
					Blocks:     1024000,
					Used:       1000,
					Available:  1023001,
					UsePercent: 1,
					MountedOn:  "/run",
				},
			},
			wantErr: false,
		},
		{
			name: "Standard output 2",
			input: `
Filesystem       1K-blocks        Used   Available Use% Mounted on
udev              32855512           0    32855512   0% /dev
tmpfs              6576452        1772     6574680   1% /run
/dev/md2        7716809892  5003753608  2324075432  69% /
/dev/md1           1012428      129976      830128  14% /boot
iac:           45097156608 31606160564 13490996044  71% /root/gamma2/iac
`,
			want: []FileSystemInfo{
				{
					Filesystem: "udev",
					Blocks:     32855512,
					Used:       0,
					Available:  32855512,
					UsePercent: 0,
					MountedOn:  "/dev",
				},
				{
					Filesystem: "tmpfs",
					Blocks:     6576452,
					Used:       1772,
					Available:  6574680,
					UsePercent: 1,
					MountedOn:  "/run",
				},
				{
					Filesystem: "/dev/md2",
					Blocks:     7716809892,
					Used:       5003753608,
					Available:  2324075432,
					UsePercent: 69,
					MountedOn:  "/",
				},
				{
					Filesystem: "/dev/md1",
					Blocks:     1012428,
					Used:       129976,
					Available:  830128,
					UsePercent: 14,
					MountedOn:  "/boot",
				},
				{
					Filesystem: "iac:",
					Blocks:     45097156608,
					Used:       31606160564,
					Available:  13490996044,
					UsePercent: 71,
					MountedOn:  "/root/gamma2/iac",
				},
			},
			wantErr: false,
		},
		{
			name: "Standard output3",
			input: `
Filesystem      1K-blocks      Used Available Use% Mounted on
u408683-sub4   1073741824 825296768 248445056  77% /home
tmpfs            1024000    1000   1023001     1% /run
`,
			want: []FileSystemInfo{
				{
					Filesystem: "u408683-sub4",
					Blocks:     1073741824,
					Used:       825296768,
					Available:  248445056,
					UsePercent: 77,
					MountedOn:  "/home",
				},
				{
					Filesystem: "tmpfs",
					Blocks:     1024000,
					Used:       1000,
					Available:  1023001,
					UsePercent: 1,
					MountedOn:  "/run",
				},
			},
			wantErr: false,
		},
		{
			name: "Empty input",
			input: `
Filesystem      1K-blocks      Used Available Use% Mounted on
`,
			want:    []FileSystemInfo{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDfOutput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDfOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDfOutput() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
