package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"storage-monitor/config"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

func setRcloneInfoMetric(version string, cfg *config.InternalConfig) {
	rcloneVersion := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rclone_info",
			Help: "Rclone executable information",
		},
		[]string{"version", "executablePath"},
	)

	prometheus.MustRegister(rcloneVersion)
	rcloneVersion.WithLabelValues(version, cfg.RcloneExecutablePath).Set(1)
}

type CommandExecutor func(name string, arg ...string) ([]byte, error)

var execCommandExecutor CommandExecutor = func(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	return cmd.Output()
}

func getRcloneVersion(rcloneExecutablePath string) (string, error) {
	return getRcloneVersionWithExecutor(rcloneExecutablePath, execCommandExecutor)
}

func getRcloneVersionWithExecutor(rcloneExecutablePath string, cmdExecutor CommandExecutor) (string, error) {
	out, err := cmdExecutor(rcloneExecutablePath, "version")

	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return "", fmt.Errorf("rclone exited with code %d", exitErr.ExitCode())
		}
		return "", fmt.Errorf("failed to execute rclone: %w", err)
	}

	firstLine := strings.SplitN(string(bytes.TrimSpace(out)), "\n", 2)[0]

	if !strings.Contains(firstLine, "rclone") {
		return "", fmt.Errorf("rclone executable is not rclone: %s", firstLine)
	}

	return firstLine, nil
}
