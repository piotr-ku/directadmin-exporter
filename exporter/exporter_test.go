package exporter

import (
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

// TestGetMetrics is a unit test for the GetMetrics function.
//
// It activates the HTTP mock, configures the response, registers
// the response function, gets the metrics using the GetMetrics function.
// The function verifies that the metrics fetched from the API match
// the expected values.
func TestGetMetrics(t *testing.T) {
	// Activate HTTP mock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Configure response
	var test = APIResponseTest{
		Response: responseFromFile("../testing/api/successful.json"),
		Status:   200,
	}

	// Register response
	httpmock.RegisterResponder("GET", fmt.Sprintf(urlFormat,
		config.Protocol, config.Username, config.Token, config.Hostname,
		config.Port), responseFunction(test))

	// Get metrics
	metrics, _ := GetMetrics(config)

	// Define tests
	tests := []struct {
		key      string
		exists   bool
		expected float64
	}{
		{key: "bandwidth", exists: true, expected: 85541},
		{key: "loadavg_five", exists: true, expected: 2.27},
		{key: "disk1", exists: false, expected: 0},
	}

	for _, test := range tests {
		metric, exists := metrics[test.key]
		assert.Equal(t, test.exists, exists)
		assert.Equal(t, test.expected, metric)
	}
}

// TestGetMetricsAPIError is a unit test for the GetMetrics function when
// the API returns an error.
//
// It activates the HTTP mock, configures the response, registers
// the response function, gets the metrics using the GetMetrics function.
// The function verifies that the metrics fetched from the API match
// the expected values.
func TestGetMetricsAPIError(t *testing.T) {
	// Activate HTTP mock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Configure response
	var test = APIResponseTest{
		Response: responseFromFile("../testing/api/invalid-token.json"),
		Status:   200,
	}

	// Register response
	httpmock.RegisterResponder("GET", fmt.Sprintf(urlFormat,
		config.Protocol, config.Username, config.Token, config.Hostname,
		config.Port), responseFunction(test))

	// Get metrics
	metrics, err := GetMetrics(config)

	// Metrics should return error
	assert.Error(t, err)
	assert.Equal(t, map[string]float64{}, metrics)
}

// TestRecordMetrics is a unit test for the RecordMetrics function.
//
// It activates the HTTP mock, configures the response, registers
// the response function, gets the metrics using the GetMetrics function.
// The function verifies that the metrics fetched from the API match
// the expected values.
func TestRecordMetrics(t *testing.T) {
	// Clean environment after test
	defer func() {
		httpmock.DeactivateAndReset()
		metrics = make(map[string]prometheus.Gauge)
	}()

	// Activate HTTP mock
	httpmock.Activate()

	// Configure response
	var test = APIResponseTest{
		Response: responseFromFile("../testing/api/successful.json"),
		Status:   200,
	}

	// Register response
	httpmock.RegisterResponder("GET", fmt.Sprintf(urlFormat,
		config.Protocol, config.Username, config.Token, config.Hostname,
		config.Port), responseFunction(test))

	// Get metrics
	metrics := RecordMetrics(config)

	// Define tests
	tests := []struct {
		key      string
		exists   bool
		expected float64
	}{
		{key: "bandwidth", exists: true},
		{key: "loadavg_five", exists: true},
		{key: "disk1", exists: false},
	}

	for _, test := range tests {
		_, exists := metrics[test.key]
		assert.Equal(t, test.exists, exists)
	}
}

// TestRecordMetricsAPIError is a unit test for the RecordMetrics function
// when the API returns an error.
//
// It activates the HTTP mock, configures the response, registers
// the response function, gets the metrics using the GetMetrics function,
// and defines the tests. The function verifies that the metrics fetched
// from the API match the expected values.
func TestRecordMetricsAPIError(t *testing.T) {
	// Clean environment after test
	defer func() {
		httpmock.DeactivateAndReset()
		metrics = make(map[string]prometheus.Gauge)
	}()

	// Activate HTTP mock
	httpmock.Activate()

	// Configure response
	var test = APIResponseTest{
		Response: responseFromFile("../testing/api/invalid-token.json"),
		Status:   200,
	}

	// Register response
	httpmock.RegisterResponder("GET", fmt.Sprintf(urlFormat,
		config.Protocol, config.Username, config.Token, config.Hostname,
		config.Port), responseFunction(test))

	// Get metrics
	metrics := RecordMetrics(config)

	// Metrics should be empty
	assert.Equal(t, map[string]prometheus.Gauge{}, metrics)
}
