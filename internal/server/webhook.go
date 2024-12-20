package server

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
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
		return "", fmt.Errorf("invalid URL path '%s'", path)
	}
	if !strings.HasPrefix(path, WebhookPath) {
		return "", fmt.Errorf("invalid URL path '%s'", path)
	}
	webhookID := path[len(WebhookPath)+1:]
	return webhookID, nil
}

// Check if the webhook id exists. Return the username
func checkWebhookID(db *sql.DB, webhookBaseURL string, webhookID string) (string, string, error) {
	list, err := getRowHashOnly(db, webhookBaseURL, webhookID)
	if err != nil || len(list) == 0 {
		return "", "", fmt.Errorf("invalid webhook ID '%s'", webhookID)
	}
	if len(list) > 1 {
		return "", "", fmt.Errorf("invalid database; found multiple webhook with webhook ID '%s'", webhookID)
	}
	return list[0].Groupname, list[0].Username, nil
}

// Read the payload from the request body
func parseWebhookPayload(req *http.Request) ([]byte, error) {
	payload, err := io.ReadAll(req.Body)
	return payload, err
}

// Write the payload to a file
func writeWebhookPayloadToFile(payloadDir string, payload []byte) error {
	payloadFilename := path.Join(payloadDir, PayLoadName)
	err := os.WriteFile(payloadFilename, payload, 0600)
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

// Process the webhook and log events
func processWebhook(c Connector, conf executeConfiguration) {
	err := ExecuteScript(c, conf)
	if err != nil {
		fmt.Printf("%s Error %s\n", time.Now().Format(time.RFC3339), err)
	}
	fmt.Printf("%s Success\n", time.Now().Format(time.RFC3339))
}

// WebhookHandler handles a HTTP POST request containing the webhook payload in its body
func (a *API) WebhookHandler(w http.ResponseWriter, req *http.Request) {
	// Check the method
	if !strings.EqualFold(req.Method, "POST") {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Error 405 - Method not allowed: invalid method: ", req.Method)
		fmt.Printf("%s Error 405 - Method not allowed: invalid method: %s\n", time.Now().Format(time.RFC3339), req.Method)
		return
	}

	// Parse and validate the request
	_, webhookID, err := parseWebhookRequest(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		fmt.Printf("%s Error %s\n", time.Now().Format(time.RFC3339), err)
		return
	}

	// Check if webhookID exists
	groupname, username, err := checkWebhookID(a.DB, a.WebhookBaseURL, webhookID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		fmt.Printf("%s Error %s\n", time.Now().Format(time.RFC3339), err)
		return
	}

	// Parse the webhook payload
	var payload []byte
	payload, err = parseWebhookPayload(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		fmt.Printf("%s Error %s\n", time.Now().Format(time.RFC3339), err)
		return
	}

	// Create the payload dir
	payloadDir := path.Join(a.DataDir, "payloads", username)
	err = os.MkdirAll(payloadDir, os.ModePerm)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		fmt.Printf("%s Error %s\n", time.Now().Format(time.RFC3339), err)
		return
	}

	// Write the payload to file
	err = writeWebhookPayloadToFile(payloadDir, payload)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Error 404 - Not found: ", err)
		fmt.Printf("%s Error %s\n", time.Now().Format(time.RFC3339), err)
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
		relayNodeName:          a.RelayNode,
		dataDir:                a.DataDir,
		homeDir:                a.HomeDir,
		webhookID:              webhookID,
		payload:                payload,
	}

	// Process the webhook in the background
	go processWebhook(a.Connector, executeConfig)

	// Succes
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Payload delivered successfully")
	fmt.Printf("%s Payload delivered successfully\n", time.Now().Format(time.RFC3339))
}
