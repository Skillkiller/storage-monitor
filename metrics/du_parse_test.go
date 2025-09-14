package metrics

import "testing"

func TestParseDuOutput(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    uint64
		wantErr bool
	}{
		{
			name:    "simple with trailing slash",
			input:   "36833627873806  gamma2/\n",
			want:    36833627873806,
			wantErr: false,
		},
		{
			name:    "tab separated path",
			input:   "4096\t/mnt/storage/\n",
			want:    4096,
			wantErr: false,
		},
		{
			name:    "spaces before long path",
			input:   "72599908815     /var/lib/docker/volumes/minecraft_aar/_data/\n",
			want:    72599908815,
			wantErr: false,
		},
		{
			name:    "second style example 1",
			input:   "19455   /var/lib/docker/volumes/minecraft_aar/_data/\n",
			want:    19455,
			wantErr: false,
		},
		{
			name:    "second style example 2",
			input:   "130816  /root/\n",
			want:    130816,
			wantErr: false,
		},
		{
			name:    "empty input returns error",
			input:   "",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDuOutput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseDuOutput() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("parseDuOutput() got = %d, want %d", got, tt.want)
			}
		})
	}
}
