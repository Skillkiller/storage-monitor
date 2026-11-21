package collectors

import (
	"storage-monitor/config"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func boolPtr(b bool) *bool {
	return &b
}

func TestRcloneCollector_WatchedDirMaxDepthZero(t *testing.T) {
	mockExecutor := func(name string, arg ...string) ([]byte, error) {
		return []byte(`[
			{"Path": "rootfile.txt", "Name": "rootfile.txt", "Size": 100, "IsDir": false},
			{"Path": "sub/nested.txt", "Name": "nested.txt", "Size": 250, "IsDir": false}
		]`), nil
	}

	collector := NewRcloneCollectorWithExecutor(
		config.RecursiveWatchedDir{
			Name:         "logs",
			MaxDepth:     0,
			Path:         "/tmp",
			FolderSize:   boolPtr(true),
			LastModified: boolPtr(false),
		},
		"rclone",
		mockExecutor,
	)

	expected := `
# HELP folder_object_count Number of objects (files and subfolders) in the given folder
# TYPE folder_object_count gauge
folder_object_count{name="logs",path="."} 2
# HELP folder_size Total size of all objects in the folder in bytes
# TYPE folder_size gauge
folder_size{name="logs",path="."} 350
# HELP watched_folder_up Whether folder metrics collection was successful (1 = OK, 0 = Error)
# TYPE watched_folder_up gauge
watched_folder_up{name="logs"} 1
`

	if err := testutil.CollectAndCompare(collector, strings.NewReader(expected)); err != nil {
		t.Fatalf("unexpected metrics output: %v", err)
	}
}

func TestRcloneCollector_RecursiveWatchedDirMaxDepthOne(t *testing.T) {
	mockExecutor := func(name string, arg ...string) ([]byte, error) {
		return []byte(`[
			{"Path": "rootfile.txt", "Name": "rootfile.txt", "Size": 100, "IsDir": false},
			{"Path": "alice/docs/file.txt", "Name": "file.txt", "Size": 400, "IsDir": false}
		]`), nil
	}

	collector := NewRcloneCollectorWithExecutor(
		config.RecursiveWatchedDir{
			Name:         "home",
			MaxDepth:     1,
			Path:         "/home",
			FolderSize:   boolPtr(true),
			LastModified: boolPtr(false),
		},
		"rclone",
		mockExecutor,
	)

	expected := `
# HELP folder_object_count Number of objects (files and subfolders) in the given folder
# TYPE folder_object_count gauge
folder_object_count{name="home",path="alice"} 1
# HELP folder_size Total size of all objects in the folder in bytes
# TYPE folder_size gauge
folder_size{name="home",path="alice"} 400
# HELP watched_folder_up Whether folder metrics collection was successful (1 = OK, 0 = Error)
# TYPE watched_folder_up gauge
watched_folder_up{name="home"} 1
`

	if err := testutil.CollectAndCompare(collector, strings.NewReader(expected)); err != nil {
		t.Fatalf("unexpected metrics output: %v", err)
	}
}

func TestRcloneCollector_InvalidJSONReturnsSingleUpZero(t *testing.T) {
	mockExecutor := func(name string, arg ...string) ([]byte, error) {
		return []byte(`invalid-json`), nil
	}

	collector := NewRcloneCollectorWithExecutor(
		config.RecursiveWatchedDir{
			Name:         "broken",
			MaxDepth:     0,
			Path:         "/broken",
			FolderSize:   boolPtr(false),
			LastModified: boolPtr(false),
		},
		"rclone",
		mockExecutor,
	)

	expected := `
# HELP watched_folder_up Whether folder metrics collection was successful (1 = OK, 0 = Error)
# TYPE watched_folder_up gauge
watched_folder_up{name="broken"} 0
`

	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(collector)

	if err := testutil.GatherAndCompare(reg, strings.NewReader(expected)); err != nil {
		t.Fatalf("unexpected metrics output: %v", err)
	}
}
