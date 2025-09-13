package main

import (
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

	if len(cfg.Hosts) == 0 {
		log.Fatal("Keine Hosts konfiguriert")
	}

	for _, host := range cfg.Hosts {
		prometheus.MustRegister(collectors.NewSshCollector(host))
	}

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
