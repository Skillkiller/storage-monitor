package collectors

import (
	"errors"
	"storage-monitor/config"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestRcloneAboutCollector_FullOutput(t *testing.T) {
	var calledArgs []string
	mockExecutor := func(name string, arg ...string) ([]byte, error) {
		calledArgs = append([]string{name}, arg...)
		return []byte(`{
			"total": 16106127360,
			"used": 3241932086,
			"trashed": 1024,
			"other": 2048,
			"free": 12864195274
		}`), nil
	}

	collector := NewRcloneAboutCollectorWithExecutor(
		config.AboutTarget{Name: "gdrive", Remote: "gdrive:"},
		"/usr/bin/rclone",
		mockExecutor,
	)

	expected := `
# HELP rclone_about_free_bytes Free storage space on the remote in bytes
# TYPE rclone_about_free_bytes gauge
rclone_about_free_bytes{name="gdrive"} 1.2864195274e+10
# HELP rclone_about_other_bytes Other storage space used on the remote in bytes
# TYPE rclone_about_other_bytes gauge
rclone_about_other_bytes{name="gdrive"} 2048
# HELP rclone_about_total_bytes Total storage capacity of the remote in bytes
# TYPE rclone_about_total_bytes gauge
rclone_about_total_bytes{name="gdrive"} 1.610612736e+10
# HELP rclone_about_trashed_bytes Trashed storage space on the remote in bytes
# TYPE rclone_about_trashed_bytes gauge
rclone_about_trashed_bytes{name="gdrive"} 1024
# HELP rclone_about_up Whether rclone about collection was successful (1 = OK, 0 = Error)
# TYPE rclone_about_up gauge
rclone_about_up{name="gdrive"} 1
# HELP rclone_about_used_bytes Used storage space on the remote in bytes
# TYPE rclone_about_used_bytes gauge
rclone_about_used_bytes{name="gdrive"} 3.241932086e+09
`

	if err := testutil.CollectAndCompare(collector, strings.NewReader(expected)); err != nil {
		t.Fatalf("unexpected metrics output: %v", err)
	}

	wantCmd := "/usr/bin/rclone about gdrive: --json"
	if gotCmd := strings.Join(calledArgs, " "); gotCmd != wantCmd {
		t.Errorf("expected command %q, got %q", wantCmd, gotCmd)
	}
}

func TestRcloneAboutCollector_PartialOutput(t *testing.T) {
	mockExecutor := func(name string, arg ...string) ([]byte, error) {
		return []byte(`{"used": 5000, "free": 10000}`), nil
	}

	collector := NewRcloneAboutCollectorWithExecutor(
		config.AboutTarget{Name: "s3", Remote: "s3:"},
		"rclone",
		mockExecutor,
	)

	expected := `
# HELP rclone_about_free_bytes Free storage space on the remote in bytes
# TYPE rclone_about_free_bytes gauge
rclone_about_free_bytes{name="s3"} 10000
# HELP rclone_about_up Whether rclone about collection was successful (1 = OK, 0 = Error)
# TYPE rclone_about_up gauge
rclone_about_up{name="s3"} 1
# HELP rclone_about_used_bytes Used storage space on the remote in bytes
# TYPE rclone_about_used_bytes gauge
rclone_about_used_bytes{name="s3"} 5000
`

	if err := testutil.CollectAndCompare(collector, strings.NewReader(expected)); err != nil {
		t.Fatalf("unexpected metrics output: %v", err)
	}
}

func TestRcloneAboutCollector_CommandError(t *testing.T) {
	mockExecutor := func(name string, arg ...string) ([]byte, error) {
		return nil, errors.New("rclone failed")
	}

	collector := NewRcloneAboutCollectorWithExecutor(
		config.AboutTarget{Name: "broken", Remote: "broken:"},
		"rclone",
		mockExecutor,
	)

	expected := `
# HELP rclone_about_up Whether rclone about collection was successful (1 = OK, 0 = Error)
# TYPE rclone_about_up gauge
rclone_about_up{name="broken"} 0
`

	if err := testutil.CollectAndCompare(collector, strings.NewReader(expected)); err != nil {
		t.Fatalf("unexpected metrics output: %v", err)
	}
}

func TestRcloneAboutCollector_InvalidJSON(t *testing.T) {
	mockExecutor := func(name string, arg ...string) ([]byte, error) {
		return []byte(`not-json`), nil
	}

	collector := NewRcloneAboutCollectorWithExecutor(
		config.AboutTarget{Name: "badjson", Remote: "badjson:"},
		"rclone",
		mockExecutor,
	)

	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(collector)

	expected := `
# HELP rclone_about_up Whether rclone about collection was successful (1 = OK, 0 = Error)
# TYPE rclone_about_up gauge
rclone_about_up{name="badjson"} 0
`

	if err := testutil.GatherAndCompare(reg, strings.NewReader(expected)); err != nil {
		t.Fatalf("unexpected metrics output: %v", err)
	}
}
