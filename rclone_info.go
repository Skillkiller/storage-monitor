package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"ssh-storage-monitor/config"
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

func getRcloneVersion(rcloneExecutablePath string) (string, error) {
	cmd := exec.Command(rcloneExecutablePath, "version")

	out, err := cmd.Output()

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
