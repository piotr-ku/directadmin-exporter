package exporter

import (
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

// TestParseResponse tests the ParseResponse function.
func TestParseResponse(t *testing.T) {
	// Define expected output
	var expected = map[string]interface{}{
		"allocated": map[string]interface{}{
			"bandwidth":   "unlimited",
			"domainptr":   "unlimited",
			"ftp":         "unlimited",
			"inode":       "unlimited",
			"mysql":       "unlimited",
			"nemailf":     "unlimited",
			"nemailml":    "unlimited",
			"nemailr":     "unlimited",
			"nemails":     "unlimited",
			"nsubdomains": "unlimited",
			"quota":       "845790",
			"vdomains":    "unlimited",
		},
		"bandwidth": "85541",
		"db_quota":  "9677621435",
		"device":    "eth0:1",
		"disk": map[string]interface{}{
			"info": map[string]interface{}{
				"columns": map[string]interface{}{
					"1024-blocks": "2",
					"Available":   "4",
					"Capacity":    "5",
					"Filesystem":  "1",
					"Mounted on":  "6",
					"Used":        "3",
				},
				"current_page": "1",
				"ipp":          "50",
				"rows":         "0",
				"total_pages":  "0",
			},
		},
		"disk1":                     "devtmpfs:7922696:0:7922696:0%:/dev",
		"disk2":                     "tmpfs:7941172:0:7941172:0%:/dev/shm",
		"disk3":                     "tmpfs:7941172:795376:7145796:11%:/run",
		"disk4":                     "tmpfs:7941172:0:7941172:0%:/sys/fs/cgroup",       // nolint: revive
		"disk5":                     "/dev/sda1:78600680:42991948:32362524:58%:/",      // nolint: revive
		"disk6":                     "/dev/sdc:46326820:24944632:19352580:57%:/volume", // nolint: revive
		"disk7":                     "/dev/sdb:655261800:589609332:38877088:94%:/home", // nolint: revive
		"disk8":                     "/dev/sda15:65390:2172:63218:4%:/boot/efi",
		"disk9":                     "tmpfs:1588232:0:1588232:0%:/run/user/0",
		"domainptr":                 "34",
		"email_deliveries":          "10917",
		"email_deliveries_incoming": "7974",
		"email_deliveries_outgoing": "2943",
		"email_quota":               "0",
		"ftp":                       "326",
		"inode":                     "6041645",
		"last_tally":                "1688682917",
		"loadavg": map[string]interface{}{
			"fifteen": "2.05",
			"five":    "2.27",
			"one":     "2.55",
		},
		"mysql":       "981",
		"nemailf":     "114",
		"nemailml":    "2",
		"nemailr":     "3",
		"nemails":     "922",
		"nresellers":  "1",
		"nsubdomains": "159",
		"nusers":      "211",
		"other_quota": "0",
		"quota":       "552070",
		"usage": map[string]interface{}{
			"bandwidth":                 "85541",
			"db_quota":                  "9677621435",
			"domainptr":                 "34",
			"email_deliveries":          "10917",
			"email_deliveries_incoming": "7974",
			"email_deliveries_outgoing": "2943",
			"email_quota":               "0",
			"ftp":                       "326",
			"inode":                     "6041645",
			"last_tally":                "1688682917",
			"mysql":                     "981",
			"nemailf":                   "114",
			"nemailml":                  "2",
			"nemailr":                   "3",
			"nemails":                   "922",
			"nresellers":                "1",
			"nsubdomains":               "159",
			"nusers":                    "211",
			"other_quota":               "0",
			"quota":                     "552070",
			"vdomains":                  "1023",
		},
		"vdomains": "1023",
	}

	// Define response
	mockResponse := APIResponseTest{
		Response: responseFromFile("../testing/api/successful.json"),
		Status:   200,
	}

	// Activate HTTP mock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register response
	httpmock.RegisterResponder("GET", fmt.Sprintf(urlFormat,
		config.Protocol, config.Username, config.Token, config.Hostname,
		config.Port), responseFunction(mockResponse))

	// Make request
	response, _ := APIRequest(config)

	// Parse request
	parsed, err := ParseResponse(response)

	// Asserts
	assert.Nil(t, err)
	assert.Equal(t, expected, parsed)
}

// TestToMetricName tests the toMetricName function.
func TestToMetricName(t *testing.T) {
	// Define tests to perform
	tests := []struct {
		given    string
		expected string
	}{
		{given: "Mounted on", expected: "mounted_on"},
		{given: "1024-blocks", expected: "1024_blocks"},
		{given: "directadmin___value", expected: "directadmin_value"},
	}

	// Perform tests
	for _, test := range tests {
		assert.Equal(t, test.expected, toMetricName(test.given))
	}
}

// TestConvertResponse tests the ConvertResponse function.
func TestConvertResponse(t *testing.T) {
	tests := []struct {
		given    map[string]interface{}
		expected map[string]float64
	}{
		{
			given: map[string]interface{}{
				"bandwidth": "85541",
				"disk1":     "devtmpfs:7922696:0:7922696:0%:/dev",
			},
			expected: map[string]float64{
				"bandwidth": 85541,
			},
		},
		{
			given: map[string]interface{}{
				"bandwidth": "85541",
				"loadavg": map[string]interface{}{
					"fifteen": "2.05",
					"five":    "2.27",
					"one":     "2.55",
				},
			},
			expected: map[string]float64{
				"bandwidth":       85541,
				"loadavg_fifteen": 2.05,
				"loadavg_five":    2.27,
				"loadavg_one":     2.55,
			},
		},
	}

	// Perform tests
	for _, test := range tests {
		assert.Equal(t, test.expected, ConvertResponse(test.given))
	}
}

// TestConvertResponseInvalidType tests the ConvertResponse function with
// an invalid type.
func TestConvertResponseInvalidType(t *testing.T) {
	// Define testing data
	given := map[string]interface{}{
		"bandwidth":  "85541",
		"unexpected": []struct{}{},
	}

	// Perform test
	assert.Panics(t, func() { ConvertResponse(given) })
}
