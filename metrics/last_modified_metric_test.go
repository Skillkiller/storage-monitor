package metrics

import "testing"

func TestParseStatOutput_ValidCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		m     uint64
		c     uint64
	}{
		{
			name:  "equal mtime and ctime #1",
			input: "1757548157;1757548157\n",
			m:     1757548157,
			c:     1757548157,
		},
		{
			name:  "equal mtime and ctime #2",
			input: "1757544094;1757544094\n",
			m:     1757544094,
			c:     1757544094,
		},
		{
			name:  "different mtime and ctime",
			input: "1747533427;1750114104\n",
			m:     1747533427,
			c:     1750114104,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, c, err := parseStatOutput(tt.input)
			if err != nil {
				t.Fatalf("parseStatOutput() unexpected error: %v", err)
			}
			if m != tt.m || c != tt.c {
				t.Fatalf("parseStatOutput() got m=%d c=%d, want m=%d c=%d", m, c, tt.m, tt.c)
			}
		})
	}
}

func TestParseStatOutput_InvalidCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "empty", input: ""},
		{name: "no separator", input: "1757548157"},
		{name: "too many fields", input: "1;2;3"},
		{name: "non numeric mtime", input: "abc;123"},
		{name: "non numeric ctime", input: "123;xyz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := parseStatOutput(tt.input)
			if err == nil {
				t.Fatalf("parseStatOutput() expected error for input %q, got none", tt.input)
			}
		})
	}
}
