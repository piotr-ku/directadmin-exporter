package exporter

import (
	"errors"
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var metrics map[string]prometheus.Gauge

// GetMetrics retrieves the metrics from the DirectAdmin API based on the
// provided configuration.
func GetMetrics(config APIConfiguration) (map[string]float64, error) {
	// Perform API Request
	response, _ := APIRequest(config)

	// Parse response
	parsed, _ := ParseResponse(response)

	// Handle API errors
	if parsed["error"] != nil {
		errMsg := parsed["error"].(string)
		log.Println(errors.New(errMsg))
		return map[string]float64{}, errors.New(errMsg)
	}

	// Convert to map[string]float64
	return ConvertResponse(parsed), nil
}

// RecordMetrics retrieves and records the metrics from the DirectAdmin API
// based on the provided configuration.
func RecordMetrics(config APIConfiguration) map[string]prometheus.Gauge {
	if metrics == nil {
		metrics = make(map[string]prometheus.Gauge)
	}
	// Get metrics
	m, err := GetMetrics(config)
	if err != nil {
		return metrics
	}

	// Return map of prometheus gauges
	for key, value := range m {
		if _, exist := metrics[key]; !exist {
			metrics[key] = promauto.NewGauge(prometheus.GaugeOpts{
				Name: "directadmin_" + key,
			})
		}
		metrics[key].Set(value)
	}

	return metrics
}
