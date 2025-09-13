package metrics

import (
	"fmt"
	"ssh-storage-monitor/config"
	ssh2 "ssh-storage-monitor/ssh"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

func GetDirSize(client *ssh.Client, dirConfig config.WatchedDir) (uint64, error) {
	result, err := ssh2.ExecuteSshCommand(client, "du -sb "+dirConfig.Path)
	if err != nil {
		return 0, fmt.Errorf("can't execute du: %w - result: %s", err, string(result))
	}

	folderSize, err := parseDuOutput(string(result))
	if err != nil {
		return 0, err
	}

	return folderSize, nil
}

func parseDuOutput(output string) (uint64, error) {
	fields := strings.Fields(output)
	if len(fields) < 1 {
		return 0, fmt.Errorf("unexpected du -sb output: %q", output)
	}
	return strconv.ParseUint(fields[0], 10, 64)
}
