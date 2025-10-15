package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"ssh-storage-monitor/collectors"
	"ssh-storage-monitor/config"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

	if len(cfg.RecursiveDirs) == 0 {
		log.Fatal("Keine watched dirs angegeben")
	}

	version, err := getRcloneVersion(cfg.RcloneExecutablePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Rclone Version: %s\n", version)

	for _, dirs := range cfg.RecursiveDirs {
		prometheus.MustRegister(collectors.NewRcloneCollector(dirs, cfg.RcloneExecutablePath))
	}

	fmt.Printf("%+v\n", cfg)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
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
