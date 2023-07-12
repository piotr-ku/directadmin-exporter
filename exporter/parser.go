package exporter

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

// toMetricName converts a metric name to a valid Prometheus metric name.
func toMetricName(name string) string {
	m := regexp.MustCompile("[^A-Za-z0-9]+")
	return strings.ToLower(m.ReplaceAllString(name, "_"))
}

// ParseResponse parses the DirectAdmin API response into a map of
// string-interface{}.
func ParseResponse(response []byte) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal(response, &data)
	return data, err
}

// ConvertResponse converts the parsed API response into
// a map of string-float64.
func ConvertResponse(response map[string]interface{}) map[string]float64 {
	data := make(map[string]float64)
	for key, value := range response {
		switch value := value.(type) {
		case string:
			float, err := strconv.ParseFloat(value, 64)
			if err == nil {
				data[toMetricName(key)] = float
			}
		case map[string]interface{}:
			for k, v := range ConvertResponse(value) {
				data[toMetricName(key+"_"+k)] = v
			}
		default:
			panic("Unexpected data type")
		}
	}
	return data
}
