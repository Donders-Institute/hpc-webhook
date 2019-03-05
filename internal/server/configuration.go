package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// ConfigurationRequest stores one row of webhook information
type ConfigurationRequest struct {
	Hash      string `json:"hash"`
	Groupname string `json:"groupname"`
	Username  string `json:"username"`
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
		return configuration, fmt.Errorf("invalid method '%s'", req.Method)
	}

	// Check the URL path
	if !isValidConfigurationURLPath(req.URL.Path) {
		return configuration, fmt.Errorf("invalid URL path '%s'", req.URL.Path)
	}

	// Obtain the configuration
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&configuration)
	if err != nil {
		return configuration, errors.New("invalid JSON body")
	}

	return configuration, err
}
