package server

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

// Webhook is an inbound github webhook
type Webhook struct {
	WebhookID string
	ID        string
	Event     string
	Signature string
	Payload   []byte
}

// Extract the webhook ID from the given url
func extractWebhookID(u *url.URL, WebhookPath string) (string, error) {
	path := u.Path
	if len(path) < len(WebhookPath) {
		return "", fmt.Errorf("Invalid URL path '%s'", path)
	}
	if !strings.HasPrefix(path, WebhookPath) {
		return "", fmt.Errorf("Invalid URL path '%s'", path)
	}
	webhookID := path[len(WebhookPath)+1:]
	return webhookID, nil
}

// Check if the webhook id exists. Return the username
func checkWebhookID(db *sql.DB, hpcWebhookHost string, hpcWebhookExternalPort string, webhookID string) (string, string, error) {
	list, err := getRowHashOnly(db, hpcWebhookHost, hpcWebhookExternalPort, webhookID)
	if err != nil || len(list) == 0 {
		return "", "", fmt.Errorf("Invalid webhook ID '%s'", webhookID)
	}
	if len(list) > 1 {
		return "", "", fmt.Errorf("Invalid database; found multiple webhook with webhook ID '%s'", webhookID)
	}
	return list[0].Groupname, list[0].Username, nil
}

// Read the payload from the request body
func parseWebhookPayload(req *http.Request) ([]byte, error) {
	payload, err := ioutil.ReadAll(req.Body)
	return payload, err
}

// Write the payload to a file
func writeWebhookPayloadToFile(payloadDir string, payload []byte, username string) error {
	payloadFilename := path.Join(payloadDir, PayLoadName)
	err := ioutil.WriteFile(payloadFilename, payload, 0600)
	if err != nil {
		return err
	}
	return nil
}

func parseWebhookRequest(req *http.Request) (*Webhook, string, error) {
	var webhook *Webhook
	var webhookID string
	var err error

	// Check the URL path
	if !isValidURLPath(req.URL.Path) {
		return webhook, "", fmt.Errorf("invalid URL path '%s'", req.URL.Path)
	}

	// Derive the webhook id (if possible)
	webhookID, err = extractWebhookID(req.URL, WebhookPath)
	if err != nil {
		return webhook, "", fmt.Errorf("invalid webhook id '%s' in URL path", webhookID)
	}
	if !isValidWebhookID(webhookID) {
		return webhook, "", fmt.Errorf("invalid webhook id '%s' in URL path", webhookID)
	}

	return webhook, webhookID, err
}

// WebhookHandler handles a HTTP POST request containing the webhook payload in its body
func (a *API) WebhookHandler(w http.ResponseWriter, req *http.Request) {
	// Check the method
	if !strings.EqualFold(req.Method, "POST") {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Printf("Error 405 - Method not allowed: invalid method: %s", req.Method)
		fmt.Fprint(w, "Error 405 - Method not allowed: invalid method: ", req.Method)
		return
	}

	// Parse and validate the request
	_, webhookID, err := parseWebhookRequest(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Check if webhookID exists
	groupname, username, err := checkWebhookID(a.DB, a.HPCWebhookHost, a.HPCWebhookExternalPort, webhookID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Parse the webhook payload
	var payload []byte
	payload, err = parseWebhookPayload(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Create the payload dir
	payloadDir := path.Join(a.DataDir, "payloads", username)
	err = os.MkdirAll(payloadDir, os.ModePerm)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Write the payload to file
	err = writeWebhookPayloadToFile(payloadDir, payload, username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Prepare the execution of the script
	payloadFilename := path.Join(payloadDir, PayLoadName)
	targetPayloadDir := path.Join(a.HomeDir, groupname, username, WebhooksWorkDir, webhookID)
	targetPayloadFilename := path.Join(targetPayloadDir, PayLoadName)
	userScriptPathFilename := path.Join(a.HomeDir, groupname, username, WebhooksWorkDir, webhookID, ScriptName)

	executeConfig := executeConfiguration{
		privateKeyFilename:     a.PrivateKeyFilename,
		payloadFilename:        payloadFilename,
		targetPayloadDir:       targetPayloadDir,
		targetPayloadFilename:  targetPayloadFilename,
		userScriptPathFilename: userScriptPathFilename,
		username:               username,
		groupname:              groupname,
		password:               a.RelayNodeTestUserPassword,
		relayNodeName:          a.RelayNode,
		dataDir:                a.DataDir,
		homeDir:                a.HomeDir,
		webhookID:              webhookID,
		payload:                payload,
	}

	// Execute the script
	err = ExecuteScript(a.Connector, executeConfig)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		return
	}

	// Succes
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Webhook handled successfully")

	return
}
