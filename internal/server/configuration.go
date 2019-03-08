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

	// Validate the configuration
	err = validateConfigurationRequest(configuration)
	if err != nil {
		return configuration, err
	}

	return configuration, err
}

// ConfigurationHandler handles a webhook registration HTTP PUT request
// with the hash and username in its body
func (a *API) ConfigurationHandler(w http.ResponseWriter, req *http.Request) {
	// Parse and validate the request
	configuration, err := parseConfigurationRequest(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Add key to authorized keys
	err = addAuthorizedPublicKey(a.HomeDir, configuration.Groupname, configuration.Username, a.PublicKeyFilename)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Add a row in the database
	err = addRow(a.DB, configuration.Hash, configuration.Groupname, configuration.Username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Succes
	webhookPayloadURL := fmt.Sprintf("https://%s:%s/webhook/%s", a.QaasHost, "443", configuration.Hash)
	configurationResponse := ConfigurationResponse{
		Webhook: webhookPayloadURL,
	}
	js, err := json.Marshal(configurationResponse)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return
}
