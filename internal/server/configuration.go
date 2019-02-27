package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// ConfigurationRequest stores one row of webhook information
type ConfigurationRequest struct {
	Hash     string `json:"hash"`
	Username string `json:"username"`
}

// ConfigurationResponse contains the complet webhook payload URL
type ConfigurationResponse struct {
	Webhook string `json:"webhook"`
}

// Parse reads and verifies the add hook in an inbound request.
func parseConfigurationRequest(req *http.Request) (ConfigurationRequest, error) {
	var configuration ConfigurationRequest
	var err error

	// Check method
	if !strings.EqualFold(req.Method, "PUT") {
		return configuration, errors.New("invalid method")
	}

	// Check the URL path
	if !isValidConfigurationURLPath(req.URL.Path) {
		return configuration, errors.New("invalid URL path")
	}

	// Obtain the configuration
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&configuration)
	if err != nil {
		return configuration, errors.New("invalid JSON body")
	}

	return configuration, err
}
