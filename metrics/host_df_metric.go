package metrics

import (
	"fmt"
	ssh2 "ssh-storage-monitor/ssh"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

func GetHostDf(client *ssh.Client) ([]FileSystemInfo, error) {
	result, err := ssh2.ExecuteSshCommand(client, "df -B 1 ")
	if err != nil {
		return nil, fmt.Errorf("can't execute df -B 1: %w - result: %s", err, string(result))
	}

	fsInfo, err := parseDfOutput(string(result))
	if err != nil {
		return nil, err
	}

	return fsInfo, nil
}

type FileSystemInfo struct {
	Filesystem string
	Blocks     uint64
	Used       uint64
	Available  uint64
	UsePercent uint8
	MountedOn  string
}

// ParseDfOutput parst einen df-Output (mehrere Zeilen) und gibt eine Liste der FileSystemInfo zurück.
func parseDfOutput(dfOutput string) ([]FileSystemInfo, error) {
	lines := strings.Split(strings.TrimSpace(dfOutput), "\n")
	var result []FileSystemInfo

	for i, line := range lines {
		// Überspringe die Header-Zeile
		if i == 0 || strings.HasPrefix(line, "Filesystem") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 6 {
			return nil, fmt.Errorf("zeile %d hat zu wenige Felder: %q", i, line)
		}
		blocks, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return nil, err
		}
		used, err := strconv.ParseUint(fields[2], 10, 64)
		if err != nil {
			return nil, err
		}
		available, err := strconv.ParseUint(fields[3], 10, 64)
		if err != nil {
			return nil, err
		}
		usePercentStr := strings.TrimSuffix(fields[4], "%")
		usePercent, err := strconv.ParseUint(usePercentStr, 10, 8)
		if err != nil {
			return nil, err
		}
		fsinfo := FileSystemInfo{
			Filesystem: fields[0],
			Blocks:     blocks,
			Used:       used,
			Available:  available,
			UsePercent: uint8(usePercent),
			MountedOn:  fields[5],
		}
		result = append(result, fsinfo)
	}
	return result, nil
}
