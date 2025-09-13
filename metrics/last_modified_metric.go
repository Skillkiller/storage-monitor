package metrics

import (
	"fmt"
	"ssh-storage-monitor/config"
	ssh2 "ssh-storage-monitor/ssh"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

func GetLastModifiedMetric(client *ssh.Client, dirConfig config.WatchedDir) (mTime uint64, cTime uint64, err error) {
	result, err := ssh2.ExecuteSshCommand(client, "stat -c '%Y;%Z' "+dirConfig.Path)
	if err != nil {
		return 0, 0, fmt.Errorf("can't execute stat: %w", err)
	}

	mTime, cTime, err = parseStatOutput(string(result))
	if err != nil {
		return 0, 0, err
	}

	return mTime, cTime, nil
}

// ParseStatOutput parst eine Ausgabe wie "1727453785;1727453785"
// und gibt die Werte als uint64 zurück.
func parseStatOutput(output string) (mTime uint64, cTime uint64, err error) {
	fields := strings.Split(strings.TrimSpace(output), ";")
	if len(fields) != 2 {
		err = fmt.Errorf("unexpected stat output: %q", output)
		return
	}

	mTime, err = strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		return
	}
	cTime, err = strconv.ParseUint(fields[1], 10, 64)
	return
}
