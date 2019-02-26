package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// Configuration stores one row of webhook information
type Configuration struct {
	Hash     string `json:"hash"`
	Username string `json:"username"`
}

// Parse reads and verifies the add hook in an inbound request.
func parseConfigurationRequest(req *http.Request) (Configuration, error) {
	var configuration Configuration
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
