package main

import (
	"log"
	"net/http"
	"storage-monitor/collectors"
	"storage-monitor/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.Println("Starting Rclone Prometheus Exporter...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if len(cfg.RecursiveDirs) == 0 {
		log.Fatal("No directories configured to watch. Exiting.")
	}

	log.Printf("Loaded configuration: %+v", cfg)

	version, err := getRcloneVersion(cfg.RcloneExecutablePath)
	if err != nil {
		log.Fatalf("Failed to get rclone version from '%s': %v", cfg.RcloneExecutablePath, err)
	}
	log.Printf("Detected rclone version: %s", version)

	setRcloneInfoMetric(version, cfg)
	log.Println("Rclone info metric registered")

	// Register collectors
	for _, dir := range cfg.RecursiveDirs {
		log.Printf("Registering collector for directory: %s", dir.Name)
		prometheus.MustRegister(collectors.NewRcloneCollector(dir, cfg.RcloneExecutablePath))
	}

	addr := ":2112"
	log.Printf("Starting metrics HTTP server on %s ...", addr)

	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
