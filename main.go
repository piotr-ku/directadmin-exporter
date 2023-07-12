package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/piotr-ku/directadmin-exporter/exporter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Define command-line flags
	port := flag.Int("port", 8080, "Port number for the HTTP server")
	ipAddress := flag.String("ip", "", "IP address for the HTTP server")
	envFile := flag.String("config", "", "Configuration file path")
	flag.Parse()

	// Get API configuration
	config := exporter.NewAPIConfiguration(*envFile)

	// Validate API configuration
	if err := exporter.ValidateAPIConfiguration(config); err != nil {
		log.Fatalln(err)
	}

	// Record metrics
	go func() {
		for {
			exporter.RecordMetrics(config)
			time.Sleep(2 * time.Second)
		}
	}()

	// Run HTTP server
	addr := fmt.Sprintf("%s:%d", *ipAddress, *port)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(addr, nil))
}
