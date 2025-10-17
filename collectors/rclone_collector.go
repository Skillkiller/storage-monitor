package collectors

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path"
	"sort"
	"ssh-storage-monitor/config"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type RcloneCollector struct {
	recursiveWatchedDirs   config.RecursiveWatchedDirs
	rcloneExecutablePath   string
	folderObjectCountDesc  *prometheus.Desc
	folderSizeDesc         *prometheus.Desc
	folderLastModifiedDesc *prometheus.Desc
}

type FileInfo struct {
	Path    string `json:"Path"`
	Name    string `json:"Name"`
	Size    int64  `json:"Size"`
	ModTime string `json:"ModTime"`
	IsDir   bool   `json:"IsDir"`
}

type FolderStats struct {
	TotalSize    int64
	ObjectCount  int
	LastModified *time.Time
}

func (c *RcloneCollector) Describe(ch chan<- *prometheus.Desc) {
	//ch <- c.folderObjectCountDesc
	//ch <- c.folderSizeDesc
	//ch <- c.folderLastModifiedDesc
}

func (c *RcloneCollector) Collect(ch chan<- prometheus.Metric) {
	args := []string{"lsjson", c.recursiveWatchedDirs.Path, "--recursive", "--no-mimetype"}

	if !*c.recursiveWatchedDirs.LastModified {
		args = append(args, "--no-modtime")
	}

	cmd := exec.Command(c.rcloneExecutablePath, args...)

	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	var files []FileInfo
	if err := json.Unmarshal(out, &files); err != nil {
		panic(err)
	}

	folderStats := make(map[string]*FolderStats)

	for _, f := range files {
		var modTime time.Time
		var parseErr error
		if !f.IsDir {
			if *c.recursiveWatchedDirs.LastModified {
				modTime, parseErr = time.Parse(time.RFC3339Nano, f.ModTime)
				if parseErr != nil {
					fmt.Printf("Warnung: Konnte Datum nicht parsen '%s': %v\n", f.ModTime, parseErr)
					continue
				}
			} else {
				continue
			}
		}

		currentPath := f.Path

		for {
			parentDir := path.Dir(currentPath)

			if parentDir == "." || parentDir == currentPath {
				break
			}

			stats, ok := folderStats[parentDir]
			if !ok {
				stats = &FolderStats{}
				folderStats[parentDir] = stats
			}

			stats.ObjectCount++

			if !f.IsDir {
				stats.TotalSize += f.Size
				if stats.LastModified == nil || modTime.After(*stats.LastModified) {
					stats.LastModified = &modTime
				}
			}

			currentPath = parentDir
		}
	}

	sortedPaths := make([]string, 0, len(folderStats))
	for path := range folderStats {
		sortedPaths = append(sortedPaths, path)
	}
	sort.Strings(sortedPaths)

	fmt.Println("--- Berechnete Ordner-Statistiken ---")
	for _, path := range sortedPaths {
		if pathDepth(path) > c.recursiveWatchedDirs.MaxDepth-1 {
			continue
		}
		stats := folderStats[path]

		lastModStr := "N/A"
		if stats.LastModified != nil {
			lastModStr = stats.LastModified.Format("2006-01-02 15:04:05")
		}

		fmt.Printf("Ordner: %-40s | Objekte: %-10d | Größe: %-15d | Letzte Änderung: %s\n", path, stats.ObjectCount, stats.TotalSize, lastModStr)

		ch <- prometheus.MustNewConstMetric(c.folderObjectCountDesc,
			prometheus.GaugeValue,
			float64(stats.ObjectCount),
			c.recursiveWatchedDirs.Name,
			path)

		if *c.recursiveWatchedDirs.FolderSize {
			ch <- prometheus.MustNewConstMetric(c.folderSizeDesc,
				prometheus.GaugeValue,
				float64(stats.TotalSize),
				c.recursiveWatchedDirs.Name,
				path)
		}

		if *c.recursiveWatchedDirs.LastModified && stats.LastModified != nil {
			ch <- prometheus.MustNewConstMetric(c.folderLastModifiedDesc,
				prometheus.GaugeValue,
				float64(stats.LastModified.Unix()),
				c.recursiveWatchedDirs.Name,
				path)
		}
	}
}

func pathDepth(pathS string) int {
	clean := path.Clean(pathS)
	return strings.Count(clean, "/")
}

func NewRcloneCollector(recursiveWatchedDirs config.RecursiveWatchedDirs, rcloneExecutablePath string) *RcloneCollector {
	return &RcloneCollector{
		recursiveWatchedDirs:   recursiveWatchedDirs,
		rcloneExecutablePath:   rcloneExecutablePath,
		folderObjectCountDesc:  prometheus.NewDesc("folder_object_count", "Number of objects (files and subfolders) in the given folder", []string{"name", "path"}, nil),
		folderSizeDesc:         prometheus.NewDesc("folder_size", "Total size of all objects in the folder in bytes", []string{"name", "path"}, nil),
		folderLastModifiedDesc: prometheus.NewDesc("folder_last_modified_time", "Unix timestamp of the most recent modification among all objects in the folder", []string{"name", "path"}, nil),
	}
}
