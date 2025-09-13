package collectors

import (
	"fmt"
	"ssh-storage-monitor/config"
	"ssh-storage-monitor/metrics"
	"ssh-storage-monitor/ssh"

	"github.com/prometheus/client_golang/prometheus"
)

type SshCollector struct {
	hostConfig config.SshHostConfig
}

func NewSshCollector(hostConfig config.SshHostConfig) *SshCollector {
	return &SshCollector{
		hostConfig: hostConfig,
	}
}

func (c *SshCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (c *SshCollector) Collect(ch chan<- prometheus.Metric) {
	fmt.Println("ssh collect")
	connection := ssh.GetSshConnection(c.hostConfig)
	defer connection.Close()

	errorCount := 0

	for _, dirConfig := range c.hostConfig.WatchedDirs {
		if *dirConfig.FolderSize {
			size, err := metrics.GetDirSize(connection, dirConfig)
			if err == nil {
				ch <- prometheus.MustNewConstMetric(
					prometheus.NewDesc(
						"folder_size",
						"Folder size in bytes",
						[]string{"host", "name"},
						nil,
					),
					prometheus.GaugeValue,
					float64(size),
					c.hostConfig.Host,
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
					prometheus.NewDesc(
						"folder_last_modified_mtime",
						"time of last data modification",
						[]string{"host", "name"},
						nil,
					),
					prometheus.GaugeValue,
					float64(mtime),
					c.hostConfig.Host,
					dirConfig.Name,
				)
				ch <- prometheus.MustNewConstMetric(
					prometheus.NewDesc(
						"folder_last_modified_ctime",
						"time of last status change",
						[]string{"host", "name"},
						nil,
					),
					prometheus.GaugeValue,
					float64(cTime),
					c.hostConfig.Host,
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
					prometheus.NewDesc(
						"folder_inode_count",
						"Number of inodes",
						[]string{"host", "name"},
						nil,
					),
					prometheus.GaugeValue,
					float64(inodeCount),
					c.hostConfig.Host,
					dirConfig.Name,
				)
			}
		}
	}

	fileSystemInfos, err := metrics.GetHostDf(connection)
	if err != nil {
		errorCount++
		fmt.Println(err)
	} else {
		for _, fsInfo := range fileSystemInfos {
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					"filesystem_total_blocks",
					"Total number of blocks in the filesystem",
					[]string{"host", "filesystem", "mount"},
					nil,
				),
				prometheus.GaugeValue,
				float64(fsInfo.Blocks),
				c.hostConfig.Host,
				fsInfo.Filesystem,
				fsInfo.MountedOn,
			)

			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					"filesystem_used_blocks",
					"Number of blocks currently used in the filesystem",
					[]string{"host", "filesystem", "mount"},
					nil,
				),
				prometheus.GaugeValue,
				float64(fsInfo.Used),
				c.hostConfig.Host,
				fsInfo.Filesystem,
				fsInfo.MountedOn,
			)

			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					"filesystem_available_blocks",
					"Number of blocks still available in the filesystem",
					[]string{"host", "filesystem", "mount"},
					nil,
				),
				prometheus.GaugeValue,
				float64(fsInfo.Available),
				c.hostConfig.Host,
				fsInfo.Filesystem,
				fsInfo.MountedOn,
			)

		}
	}

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"error_status",
			"Indicates if there was an error while collecting metrics",
			[]string{"host"},
			nil,
		),
		prometheus.GaugeValue,
		float64(errorCount),
		c.hostConfig.Host)
}
