package collectors

import (
	"encoding/json"
	"fmt"
	"os/exec"
	pathlib "path"
	"ssh-storage-monitor/config"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type RcloneCollector struct {
	recursiveWatchedDir    config.RecursiveWatchedDir
	rcloneExecutablePath   string
	folderObjectCountDesc  *prometheus.Desc
	folderSizeDesc         *prometheus.Desc
	folderLastModifiedDesc *prometheus.Desc
	watchedFolderUp        *prometheus.Desc
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
	args := []string{"lsjson", c.recursiveWatchedDir.Path, "--recursive", "--no-mimetype"}

	if !*c.recursiveWatchedDir.LastModified {
		args = append(args, "--no-modtime")
	}

	cmd := exec.Command(c.rcloneExecutablePath, args...)

	out, err := cmd.Output()
	if err != nil {
		sendError(c, ch)
		return
	}

	var files []FileInfo
	if err := json.Unmarshal(out, &files); err != nil {
		sendError(c, ch)
	}

	folderStats := make(map[string]*FolderStats)

	for _, f := range files {
		var modTime time.Time
		var parseErr error
		if !f.IsDir {
			if *c.recursiveWatchedDir.LastModified {
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
			parentDir := pathlib.Dir(currentPath)

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

	for path := range folderStats {
		if pathDepth(path) > c.recursiveWatchedDir.MaxDepth-1 {
			continue
		}
		stats := folderStats[path]

		ch <- prometheus.MustNewConstMetric(c.folderObjectCountDesc,
			prometheus.GaugeValue,
			float64(stats.ObjectCount),
			c.recursiveWatchedDir.Name,
			path)

		if *c.recursiveWatchedDir.FolderSize {
			ch <- prometheus.MustNewConstMetric(c.folderSizeDesc,
				prometheus.GaugeValue,
				float64(stats.TotalSize),
				c.recursiveWatchedDir.Name,
				path)
		}

		if *c.recursiveWatchedDir.LastModified && stats.LastModified != nil {
			ch <- prometheus.MustNewConstMetric(c.folderLastModifiedDesc,
				prometheus.GaugeValue,
				float64(stats.LastModified.Unix()),
				c.recursiveWatchedDir.Name,
				path)
		}
	}
	ch <- prometheus.MustNewConstMetric(c.watchedFolderUp, prometheus.GaugeValue, 1, c.recursiveWatchedDir.Name)
}

func sendError(c *RcloneCollector, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.watchedFolderUp, prometheus.GaugeValue, 0, c.recursiveWatchedDir.Name)
}

func pathDepth(pathS string) int {
	clean := pathlib.Clean(pathS)
	return strings.Count(clean, "/")
}

func NewRcloneCollector(recursiveWatchedDirs config.RecursiveWatchedDir, rcloneExecutablePath string) *RcloneCollector {
	return &RcloneCollector{
		recursiveWatchedDir:    recursiveWatchedDirs,
		rcloneExecutablePath:   rcloneExecutablePath,
		folderObjectCountDesc:  prometheus.NewDesc("folder_object_count", "Number of objects (files and subfolders) in the given folder", []string{"name", "path"}, nil),
		folderSizeDesc:         prometheus.NewDesc("folder_size", "Total size of all objects in the folder in bytes", []string{"name", "path"}, nil),
		folderLastModifiedDesc: prometheus.NewDesc("folder_last_modified_time", "Unix timestamp of the most recent modification among all objects in the folder", []string{"name", "path"}, nil),
		watchedFolderUp:        prometheus.NewDesc("watched_folder_up", "Whether folder metrics collection was successful (1 = OK, 0 = Error)", []string{"name"}, nil),
	}
}
