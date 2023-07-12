package exporter

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// API configuration
var config = APIConfiguration{
	Hostname: "localhost",
	Protocol: "http",
	Port:     "2222",
	Username: "admin",
	Token:    "SECRET",
}

// unsetEnvironment unsets the environment variables specified in the filenames.
func unsetEnvironment(filenames ...string) {
	variables, err := godotenv.Read(filenames...)
	if err != nil {
		panic(err)
	}
	for variable := range variables {
		err := os.Unsetenv(variable)
		if err != nil {
			panic(err)
		}
	}
}

// TestNewAPIConfiguration tests the NewAPIConfiguration function.
func TestNewAPIConfiguration(t *testing.T) {
	// Get environment files to load
	filenames := filepath.Join("..", ".env.test")

	// Create APIConfiguration struct
	given := NewAPIConfiguration(filenames)

	// Unset environment variables after the test
	defer unsetEnvironment(filenames)

	// Expected result
	expected := APIConfiguration{
		Hostname: "localhost",
		Protocol: "http",
		Port:     "2222",
		Username: "admin",
		Token:    "SECRET_TOKEN",
	}

	// Test
	assert.Equal(t, expected, given)
}

// TestNewAPIConfigurationNonExistingFile tests the NewAPIConfiguration function
// with a non-existing file.
func TestNewAPIConfigurationNonExistingFile(t *testing.T) {
	// Get environment files to load
	filenames := filepath.Join("..", ".env.non-existing-file")

	// Create APIConfiguration struct
	given := NewAPIConfiguration(filenames)

	// Expected result
	expected := APIConfiguration{}

	// Test
	assert.Equal(t, expected, given)
}

// TestValidateAPIConfiguration tests the ValidateAPIConfiguration function.
func TestValidateAPIConfiguration(t *testing.T) {
	// Define tests
	tests := []struct {
		name     string
		config   APIConfiguration
		expected error
	}{
		{
			name: "Valid configuration",
			config: APIConfiguration{
				Hostname: "s1.hostname.com",
				Protocol: "http",
				Port:     "2222",
				Username: "admin",
				Token:    "SECRET",
			},
			expected: nil,
		},
		{
			name: "Valid configuration with IP address as a host",
			config: APIConfiguration{
				Hostname: "127.0.0.1",
				Protocol: "http",
				Port:     "2222",
				Username: "admin",
				Token:    "SECRET",
			},
			expected: nil,
		},
		{
			name: "Valid configuration with localhost as a host",
			config: APIConfiguration{
				Hostname: "localhost",
				Protocol: "http",
				Port:     "2222",
				Username: "admin",
				Token:    "SECRET",
			},
			expected: nil,
		},
		{
			name: "Missing hostname",
			config: APIConfiguration{
				Hostname: "",
				Protocol: "http",
				Port:     "2222",
				Username: "admin",
				Token:    "SECRET",
			},
			expected: errors.New("Missing hostname"),
		},
		{
			name: "Invalid hostname",
			config: APIConfiguration{
				Hostname: "---",
				Protocol: "http",
				Port:     "2222",
				Username: "admin",
				Token:    "SECRET",
			},
			expected: errors.New("Invalid hostname"),
		},
		{
			name: "Missing protocol",
			config: APIConfiguration{
				Hostname: "s1.hostname.com",
				Protocol: "",
				Port:     "2222",
				Username: "admin",
				Token:    "SECRET",
			},
			expected: errors.New("Missing protocol"),
		},
		{
			name: "Invalid protocol",
			config: APIConfiguration{
				Hostname: "s1.hostname.com",
				Protocol: "xxx",
				Port:     "2222",
				Username: "admin",
				Token:    "SECRET",
			},
			expected: errors.New("Invalid protocol"),
		},
		{
			name: "Missing port",
			config: APIConfiguration{
				Hostname: "s1.hostname.com",
				Protocol: "http",
				Port:     "",
				Username: "admin",
				Token:    "SECRET",
			},
			expected: errors.New("Missing port"),
		},
		{
			name: "Invalid port",
			config: APIConfiguration{
				Hostname: "s1.hostname.com",
				Protocol: "http",
				Port:     "invalid",
				Username: "admin",
				Token:    "SECRET",
			},
			expected: errors.New("Invalid port"),
		},
		{
			name: "Missing username",
			config: APIConfiguration{
				Hostname: "s1.hostname.com",
				Protocol: "http",
				Port:     "2222",
				Username: "",
				Token:    "SECRET",
			},
			expected: errors.New("Missing username"),
		},
		{
			name: "Missing token",
			config: APIConfiguration{
				Hostname: "s1.hostname.com",
				Protocol: "http",
				Port:     "2222",
				Username: "admin",
				Token:    "",
			},
			expected: errors.New("Missing token"),
		},
	}

	// Run tests
	for _, test := range tests {
		// Make validation
		err := ValidateAPIConfiguration(test.config)
		if test.expected == nil {
			assert.Nil(t, err, "%s: %v", test.name, err)
		} else {
			assert.Error(t, err, "%s: %v", test.name, err)
		}
	}
}

// APIResponseTest represents a test case for APIRequest function.
type APIResponseTest struct {
	Response string
	Status   int
}

// responseFunction is a helper function to create a response function for
// HTTP mocking.
func responseFunction(test APIResponseTest) func(*http.Request) (*http.Response,
	error) {
	// Prepare response
	return func(req *http.Request) (*http.Response, error) {
		if test.Status >= 400 {
			return httpmock.NewStringResponse(test.Status,
				""), fmt.Errorf("Error %d", test.Status)
		}

		return httpmock.NewStringResponse(test.Status,
			test.Response), nil
	}
}

// responseFromFile returns the content of a file as a string for HTTP mocking.
func responseFromFile(file string) string {
	return httpmock.File(file).String()
}

// TestAPIRequest tests the APIRequest function.
func TestAPIRequest(t *testing.T) {
	// Tests to perform
	tests := []APIResponseTest{
		{
			Response: responseFromFile("../testing/api/successful.json"),
			Status:   200,
		},
		{
			Response: responseFromFile("../testing/api/invalid-token.json"),
			Status:   200,
		},
		{
			Response: "",
			Status:   403,
		},
		{
			Response: "",
			Status:   404,
		},
		{
			Response: "",
			Status:   500,
		},
	}

	for _, test := range tests {
		// Activate HTTP mock
		httpmock.Activate()

		// Register response
		httpmock.RegisterResponder("GET", fmt.Sprintf(urlFormat,
			config.Protocol, config.Username, config.Token, config.Hostname,
			config.Port), responseFunction(test))

		// Make request
		response, err := APIRequest(config)

		// Check error
		if test.Status >= 400 {
			assert.Error(t, err)
		} else {
			assert.Nil(t, err)
		}

		// Check response
		assert.Equal(t, test.Response, bytes.NewBuffer(response).String())

		// Deactivate and reset HTTP mock
		httpmock.DeactivateAndReset()
	}
}

// TestAPIRequestIOUtilReadAllError tests the APIRequest function when
// ioutil.ReadAll returns an error.
func TestAPIRequestIOUtilReadAllError(t *testing.T) {
	// mock ioutil.ReadAll()
	mockIOReadAll = func(r io.Reader) ([]byte, error) {
		return []byte{}, errors.New("faked ioutil.ReadAll() error")
	}
	defer func() {
		mockIOReadAll = io.ReadAll
	}()

	// Activate HTTP mock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Configure response
	var test = APIResponseTest{
		Response: "../testing/api/successfull.json",
		Status:   200,
	}

	// Register response
	httpmock.RegisterResponder("GET", fmt.Sprintf(urlFormat,
		config.Protocol, config.Username, config.Token, config.Hostname,
		config.Port), responseFunction(test))

	// Make request
	response, err := APIRequest(config)

	// Check error
	assert.Error(t, err)

	// Check response
	assert.Equal(t, "", bytes.NewBuffer(response).String())
}
