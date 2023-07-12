package exporter

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

var urlFormat = "%s://%s:%s@%s:%s/CMD_API_ADMIN_STATS?json=yes"
var mockIOReadAll = io.ReadAll

// APIConfiguration represents the configuration data for the API.
type APIConfiguration struct {
	Hostname string `validate:"required,hostname|ip"`
	Protocol string `validate:"required,oneof=http https"`
	Port     string `validate:"required,number"`
	Username string `validate:"required"`
	Token    string `validate:"required"`
}

// NewAPIConfiguration returns a new APIConfiguration struct filled with data
// from the environment variables.
func NewAPIConfiguration(filenames ...string) APIConfiguration {
	// Get environment variables
	err := godotenv.Load(filenames...)
	if err != nil {
		log.Println(err)
	}

	return APIConfiguration{
		Hostname: os.Getenv("DIRECTADMIN_HOSTNAME"),
		Protocol: os.Getenv("DIRECTADMIN_PROTOCOL"),
		Port:     os.Getenv("DIRECTADMIN_PORT"),
		Username: os.Getenv("DIRECTADMIN_USERNAME"),
		Token:    os.Getenv("DIRECTADMIN_TOKEN"),
	}
}

// ValidateAPIConfiguration validates the APIConfiguration data.
func ValidateAPIConfiguration(config APIConfiguration) error {
	validate := validator.New()
	return validate.Struct(config)
}

// APIRequest performs a request to the DirectAdmin API.
func APIRequest(config APIConfiguration) ([]byte, error) {
	// Perform a request to the DirectAdmin API
	resp, err := http.Get(fmt.Sprintf(urlFormat, config.Protocol,
		config.Username, config.Token, config.Hostname, config.Port))
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := mockIOReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	// Return the response body
	return body, nil
}
