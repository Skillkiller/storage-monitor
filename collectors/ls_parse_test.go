package collectors

import (
	"reflect"
	"testing"
)

func TestParseLsOutput_Examples(t *testing.T) {
	base := "/base"

	cases := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name: "example block 1 (mix of dirs and files)",
			input: "total 7881880\n" +
				"drwxr-xr-x   1 root   root            0 Sep 11 15:41 iac\n" +
				"-rw-------   1 root   root        44306 Sep 10 23:35 rclone2017.config\n" +
				"drwxr-xr-x 234 root   root        12288 Sep 10 23:25 v2\n" +
				"-rw-r--r--   1 ubuntu ubuntu 8070983680 Sep 10 23:43 v2.tar\n",
			want: []string{base + "/iac", base + "/v2"},
		},
		{
			name: "example block 2 (only dirs)",
			input: "total 12\n" +
				"drwx--x--x 4 root root 4096 Sep 24  2024 containerd\n" +
				"drwxr-xr-x 7 root root 4096 Sep 25  2024 docker\n" +
				"drwxr-xr-x 2 root root 4096 Sep 27  2024 logs\n",
			want: []string{base + "/containerd", base + "/docker", base + "/logs"},
		},
		{
			name: "example block 3 (single dir)",
			input: "total 4\n" +
				"drwxr-xr-x 2 root root 4096 Sep  2 00:18 storage\n",
			want: []string{base + "/storage"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseLsOutput(tc.input, base)
			if !reflect.DeepEqual(got, tc.want) {
				// Use %#v to show slice content clearly
				t.Fatalf("parseLsOutput() got = %#v, want %#v", got, tc.want)
			}
		})
	}
}
