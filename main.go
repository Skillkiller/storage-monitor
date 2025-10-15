package main

import (
	"fmt"
	"log"
	"net/http"
	"ssh-storage-monitor/collectors"
	"ssh-storage-monitor/config"

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

	setRcloneInfoMetric(version, cfg)

	for _, dirs := range cfg.RecursiveDirs {
		prometheus.MustRegister(collectors.NewRcloneCollector(dirs, cfg.RcloneExecutablePath))
	}

	fmt.Printf("%+v\n", cfg)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
