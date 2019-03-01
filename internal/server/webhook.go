package server

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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
	webhookID := path[len(WebhookPath):]
	return webhookID, nil
}

// Check if the webhook id exists. Return the unsername
func checkWebhookID(db *sql.DB, webhookID string) (string, error) {
	list, err := getRow(db, webhookID)
	if err != nil || len(list) == 0 {
		return "", fmt.Errorf("Invalid webhook ID '%s'", webhookID)
	}
	if len(list) > 1 {
		return "", fmt.Errorf("Invalid database; found multiple webhook with webhook ID '%s'", webhookID)
	}
	return list[0].Username, nil
}

func parseWebhookPayload(req *http.Request) ([]byte, error) {
	payload, err := ioutil.ReadAll(req.Body)
	return payload, err
}

func parseWebhookRequest(req *http.Request) (*Webhook, string, error) {
	var webhook *Webhook
	var webhookID string
	var err error

	// Check the method
	if !strings.EqualFold(req.Method, "POST") {
		return webhook, "", errors.New("invalid method")
	}

	// Check the URL path
	if !isValidURLPath(req.URL.Path) {
		return webhook, "", errors.New("invalid URL path")
	}

	// Derive the webhook id (if possible)
	webhookID, err = extractWebhookID(req.URL, WebhookPath)
	if err != nil {
		return webhook, "", errors.New("invalid webhook id in URL path")
	}
	if !isValidWebhookID(webhookID) {
		return webhook, "", errors.New("invalid webhook id in URL path")
	}

	return webhook, webhookID, err
}
