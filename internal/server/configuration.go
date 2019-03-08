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

// ConfigurationResponse contains the complete webhook payload URL
type ConfigurationResponse struct {
	Webhook string `json:"webhook"`
}

// ConfigurationInfoResponse contains the detailed information about a specific webhook
type ConfigurationInfoResponse struct {
	Webhook Item `json:"webhook"`
}

// ConfigurationListResponse contains the list of regstered webhooks for a certain user
type ConfigurationListResponse struct {
	Webhooks []Item `json:"webhooks"`
}

// ConfigurationDeleteResponse contains the webhook that has been deleted
type ConfigurationDeleteResponse struct {
	Webhook string `json:"webhook"`
}

func parseConfigurationAddRequest(req *http.Request) (ConfigurationRequest, error) {
	var configuration ConfigurationRequest
	var err error

	// Check method
	if !strings.EqualFold(req.Method, "PUT") {
		return configuration, fmt.Errorf("invalid method '%s'", req.Method)
	}

	// Check the URL path
	if !isValidConfigurationAddURLPath(req.URL.Path) {
		return configuration, fmt.Errorf("invalid URL path '%s'", req.URL.Path)
	}

	// Obtain the configuration
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&configuration)
	if err != nil {
		return configuration, errors.New("invalid JSON body")
	}

	// Validate the configuration
	validateHash := true
	err = validateConfigurationRequest(configuration, validateHash)
	if err != nil {
		return configuration, err
	}

	return configuration, err
}

func parseConfigurationInfoRequest(req *http.Request) (ConfigurationRequest, error) {
	var configuration ConfigurationRequest
	var err error

	// Check method
	if !strings.EqualFold(req.Method, "GET") {
		return configuration, fmt.Errorf("invalid method '%s'", req.Method)
	}

	// Check the URL path
	if !isValidConfigurationInfoURLPath(req.URL.Path) {
		return configuration, fmt.Errorf("invalid URL path '%s'", req.URL.Path)
	}

	// Obtain the configuration
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&configuration)
	if err != nil {
		return configuration, errors.New("invalid JSON body")
	}

	// Validate the configuration
	validateHash := true
	err = validateConfigurationRequest(configuration, validateHash)
	if err != nil {
		return configuration, err
	}

	return configuration, err
}

func parseConfigurationListRequest(req *http.Request) (ConfigurationRequest, error) {
	var configuration ConfigurationRequest
	var err error

	// Check method
	if !strings.EqualFold(req.Method, "GET") {
		return configuration, fmt.Errorf("invalid method '%s'", req.Method)
	}

	// Check the URL path
	if !isValidConfigurationListURLPath(req.URL.Path) {
		return configuration, fmt.Errorf("invalid URL path '%s'", req.URL.Path)
	}

	// Obtain the configuration
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&configuration)
	if err != nil {
		return configuration, errors.New("invalid JSON body")
	}

	// Validate the configuration
	validateHash := false
	err = validateConfigurationRequest(configuration, validateHash)
	if err != nil {
		return configuration, err
	}

	return configuration, err
}

func parseConfigurationDeleteRequest(req *http.Request) (ConfigurationRequest, error) {
	var configuration ConfigurationRequest
	var err error

	// Check method
	if !strings.EqualFold(req.Method, "DELETE") {
		return configuration, fmt.Errorf("invalid method '%s'", req.Method)
	}

	// Check the URL path
	if !isValidConfigurationDeleteURLPath(req.URL.Path) {
		return configuration, fmt.Errorf("invalid URL path '%s'", req.URL.Path)
	}

	// Obtain the configuration
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&configuration)
	if err != nil {
		return configuration, errors.New("invalid JSON body")
	}

	// Validate the configuration
	validateHash := true
	err = validateConfigurationRequest(configuration, validateHash)
	if err != nil {
		return configuration, err
	}

	return configuration, err
}

// ConfigurationAddHandler handles a HTTP PUT request
// to register a certain webhook with hash, groupname, and username in its body
func (a *API) ConfigurationAddHandler(w http.ResponseWriter, req *http.Request) {
	// Parse and validate the request
	configuration, err := parseConfigurationAddRequest(req)
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

// ConfigurationInfoHandler handles a HTTP GET request
// to obtain detailed information about a specific webhook
func (a *API) ConfigurationInfoHandler(w http.ResponseWriter, req *http.Request) {
	// Parse and validate the request
	configuration, err := parseConfigurationInfoRequest(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Get the item
	list, err := getRow(a.DB, configuration.Hash)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}
	item := list[0]

	// Succes
	configurationInfoResponse := ConfigurationInfoResponse{
		Webhook: item,
	}
	js, err := json.Marshal(configurationInfoResponse)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	fmt.Printf("%+v\n", string(js))
	return
}

// ConfigurationListHandler handles a HTTP GET request
// to obtain all webhooks for a certain user
func (a *API) ConfigurationListHandler(w http.ResponseWriter, req *http.Request) {
	// Parse and validate the request
	configuration, err := parseConfigurationListRequest(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Get the list of webhooks
	list, err := getListRows(a.DB, configuration.Groupname, configuration.Username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Succes
	configurationListResponse := ConfigurationListResponse{
		Webhooks: list,
	}
	js, err := json.Marshal(configurationListResponse)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	fmt.Printf("%+v\n", string(js))
	return
}

// ConfigurationDeleteHandler handles a HTTP DELETE request
// to delete a certain webhook for a certain user
func (a *API) ConfigurationDeleteHandler(w http.ResponseWriter, req *http.Request) {
	// Parse and validate the request
	configuration, err := parseConfigurationDeleteRequest(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Delete a row in the database
	err = deleteRow(a.DB, configuration.Hash, configuration.Groupname, configuration.Username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Succes
	configurationDeleteResponse := ConfigurationDeleteResponse{
		Webhook: configuration.Hash,
	}
	js, err := json.Marshal(configurationDeleteResponse)
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
