package main

import (
	"fmt"
	"log"
	"net/http"
	"ssh-storage-monitor/config"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

	if len(cfg.RecursiveDirs)+len(cfg.WatchedDirs) == 0 {
		log.Fatal("Keine watched dirs angegeben")
	}

	//for _, host := range cfg.Hosts {
	//	prometheus.MustRegister(collectors.NewSshCollector(host))
	//}

	fmt.Printf("%+v\n", cfg)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
