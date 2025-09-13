package metrics

import (
	"fmt"
	"ssh-storage-monitor/config"
	ssh2 "ssh-storage-monitor/ssh"

	"golang.org/x/crypto/ssh"
)

func GetInodeCount(client *ssh.Client, dirConfig config.WatchedDir) (uint64, error) {
	result, err := ssh2.ExecuteSshCommand(client, "du -s --inodes "+dirConfig.Path)
	if err != nil {
		return 0, fmt.Errorf("can't execute du -s --inodes: %w", err)
	}

	folderSize, err := parseDuOutput(string(result))
	if err != nil {
		return 0, err
	}

	return folderSize, nil
}
