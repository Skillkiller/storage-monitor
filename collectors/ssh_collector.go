package collectors

import (
	"fmt"
	"ssh-storage-monitor/config"
	"ssh-storage-monitor/metrics"
	"ssh-storage-monitor/ssh"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	ssh2 "golang.org/x/crypto/ssh"
)

type SshCollector struct {
	hostConfig            config.SshHostConfig
	folderSizeDesc        *prometheus.Desc
	lastModifiedMtimeDesc *prometheus.Desc
	lastModifiedCtimeDesc *prometheus.Desc
	inodeCountDesc        *prometheus.Desc
	fsTotalBlocksDesc     *prometheus.Desc
	fsUsedBlocksDesc      *prometheus.Desc
	fsAvailableBlocksDesc *prometheus.Desc
	errorStatusDesc       *prometheus.Desc
}

func NewSshCollector(hostConfig config.SshHostConfig) *SshCollector {
	return &SshCollector{
		hostConfig:            hostConfig,
		folderSizeDesc:        prometheus.NewDesc("folder_size", "Folder size in bytes", []string{"target_name", "folder_name"}, nil),
		lastModifiedMtimeDesc: prometheus.NewDesc("folder_last_modified_mtime", "time of last data modification", []string{"target_name", "folder_name"}, nil),
		lastModifiedCtimeDesc: prometheus.NewDesc("folder_last_modified_ctime", "time of last status change", []string{"target_name", "folder_name"}, nil),
		inodeCountDesc:        prometheus.NewDesc("folder_inode_count", "Number of inodes", []string{"target_name", "folder_name"}, nil),
		fsTotalBlocksDesc:     prometheus.NewDesc("filesystem_total_blocks", "Total number of blocks in the filesystem", []string{"target_name", "filesystem", "mount"}, nil),
		fsUsedBlocksDesc:      prometheus.NewDesc("filesystem_used_blocks", "Number of blocks currently used in the filesystem", []string{"target_name", "filesystem", "mount"}, nil),
		fsAvailableBlocksDesc: prometheus.NewDesc("filesystem_available_blocks", "Number of blocks still available in the filesystem", []string{"target_name", "filesystem", "mount"}, nil),
		errorStatusDesc:       prometheus.NewDesc("error_status", "Indicates if there was an error while collecting metrics", []string{"target_name"}, nil),
	}
}

func (c *SshCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.folderSizeDesc
	ch <- c.lastModifiedMtimeDesc
	ch <- c.lastModifiedCtimeDesc
	ch <- c.inodeCountDesc
	ch <- c.fsTotalBlocksDesc
	ch <- c.fsUsedBlocksDesc
	ch <- c.fsAvailableBlocksDesc
	ch <- c.errorStatusDesc
}

func (c *SshCollector) Collect(ch chan<- prometheus.Metric) {
	fmt.Println("ssh collect")
	connection := ssh.GetSshConnection(c.hostConfig)
	defer connection.Close()

	errorCount := 0

	recursiveDirs := []config.WatchedDir{}

	for _, dirConfig := range c.hostConfig.RecursiveDirs {
		directories := getRecursiveDirectories(connection, dirConfig)
		recursiveDirs = append(recursiveDirs, directories...)
	}

	allWatchedDirs := append(recursiveDirs, c.hostConfig.WatchedDirs...)

	for _, dirConfig := range allWatchedDirs {
		if *dirConfig.FolderSize {
			size, err := metrics.GetDirSize(connection, dirConfig)
			if err == nil {
				ch <- prometheus.MustNewConstMetric(
					c.folderSizeDesc,
					prometheus.GaugeValue,
					float64(size),
					c.hostConfig.Name,
					dirConfig.Name)
			} else {
				errorCount++
				fmt.Println(err)
			}
		}

		if *dirConfig.LastModified {
			mtime, cTime, err := metrics.GetLastModifiedMetric(connection, dirConfig)
			if err != nil {
				errorCount++
				fmt.Println(err)
			} else {
				ch <- prometheus.MustNewConstMetric(
					c.lastModifiedMtimeDesc,
					prometheus.GaugeValue,
					float64(mtime),
					c.hostConfig.Name,
					dirConfig.Name,
				)
				ch <- prometheus.MustNewConstMetric(
					c.lastModifiedCtimeDesc,
					prometheus.GaugeValue,
					float64(cTime),
					c.hostConfig.Name,
					dirConfig.Name,
				)

			}
		}

		if *dirConfig.InodeCount {
			inodeCount, err := metrics.GetInodeCount(connection, dirConfig)
			if err != nil {
				errorCount++
				fmt.Println(err)
			} else {
				ch <- prometheus.MustNewConstMetric(
					c.inodeCountDesc,
					prometheus.GaugeValue,
					float64(inodeCount),
					c.hostConfig.Name,
					dirConfig.Name,
				)
			}
		}
	}

	if *c.hostConfig.FilesystemMetrics {
		fileSystemInfos, err := metrics.GetHostDf(connection)
		if err != nil {
			errorCount++
			fmt.Println(err)
		} else {
			for _, fsInfo := range fileSystemInfos {
				ch <- prometheus.MustNewConstMetric(
					c.fsTotalBlocksDesc,
					prometheus.GaugeValue,
					float64(fsInfo.Blocks),
					c.hostConfig.Name,
					fsInfo.Filesystem,
					fsInfo.MountedOn,
				)

				ch <- prometheus.MustNewConstMetric(
					c.fsUsedBlocksDesc,
					prometheus.GaugeValue,
					float64(fsInfo.Used),
					c.hostConfig.Name,
					fsInfo.Filesystem,
					fsInfo.MountedOn,
				)

				ch <- prometheus.MustNewConstMetric(
					c.fsAvailableBlocksDesc,
					prometheus.GaugeValue,
					float64(fsInfo.Available),
					c.hostConfig.Name,
					fsInfo.Filesystem,
					fsInfo.MountedOn,
				)

			}
		}
	}

	ch <- prometheus.MustNewConstMetric(
		c.errorStatusDesc,
		prometheus.GaugeValue,
		float64(errorCount),
		c.hostConfig.Name)
}

func getRecursiveDirectories(client *ssh2.Client, dirConfig config.RecursiveWatchedDirs) []config.WatchedDir {
	maxDepth := dirConfig.MaxDepth
	allDirs := []string{dirConfig.Path}
	queue := []string{dirConfig.Path}

	for depth := 1; depth <= maxDepth && len(queue) > 0; depth++ {
		var next []string
		for _, dir := range queue {
			dirs, err := listDirs(client, dir)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			allDirs = append(allDirs, dirs...)
			next = append(next, dirs...)
		}
		queue = next
	}

	var watchedDirs []config.WatchedDir
	for _, dir := range allDirs {
		watchedDirs = append(watchedDirs, config.WatchedDir{
			Name:         dir,
			Path:         dir,
			InodeCount:   dirConfig.InodeCount,
			FolderSize:   dirConfig.FolderSize,
			LastModified: dirConfig.LastModified,
		})
	}
	return watchedDirs
}

func listDirs(client *ssh2.Client, path string) ([]string, error) {
	output, err := ssh.ExecuteSshCommand(client, "ls -l "+path)
	if err != nil {
		return nil, err
	}

	var dirs []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}
		// Zeilen von `ls -l` die mit "d" anfangen sind Verzeichnisse
		if strings.HasPrefix(fields[0], "d") {
			dirs = append(dirs, path+"/"+fields[len(fields)-1])
		}
	}
	return dirs, nil
}
