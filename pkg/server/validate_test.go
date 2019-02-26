package server

import (
	"testing"
)

func TestValidConfigurationURLPath(t *testing.T) {
	cases := []struct {
		urlPath        string
		expectedResult bool
	}{
		{
			urlPath:        "/configuration",
			expectedResult: true, // Valid configuration URL path, no error
		},
		{
			urlPath:        "/configuration/",
			expectedResult: false, // Invalid configuration URL path
		},
		{
			urlPath:        "/configuration/nonexisting",
			expectedResult: false, // Invalid configuration URL path
		},
		{
			urlPath:        "/nonexisting/",
			expectedResult: false, // Invalid configuration URL path
		},
	}

	for _, c := range cases {
		result := isValidConfigurationURLPath(c.urlPath)
		if result != c.expectedResult {
			if c.expectedResult {
				t.Errorf("Expected valid url path '%s', but got invalid url path", c.urlPath)
			} else {
				t.Errorf("Expected invalid url path '%s', but got valid url path", c.urlPath)
			}
		}
	}
}

func TestValidURLPath(t *testing.T) {
	cases := []struct {
		urlPath        string
		expectedResult bool
	}{
		{
			urlPath:        "/webhook/e66d248b67c0442fe2cbad7e248651fd4569ee8ecc72ee5a19b0e55ac1ef4492",
			expectedResult: true, // Valid url path, no error
		},
		{
			urlPath:        "/webhook/",
			expectedResult: false, // Invalid url path
		},
		{
			urlPath:        "/webhook",
			expectedResult: false, // Invalid url path
		},
		{
			urlPath:        "/nonexisting/e66d248b67c0442fe2cbad7e248651fd4569ee8ecc72ee5a19b0e55ac1ef4492",
			expectedResult: false, // Invalid url path
		},
	}

	for _, c := range cases {
		result := isValidURLPath(c.urlPath)
		if result != c.expectedResult {
			if c.expectedResult {
				t.Errorf("Expected valid url path '%s', but got invalid url path", c.urlPath)
			} else {
				t.Errorf("Expected invalid url path '%s', but got valid url path", c.urlPath)
			}
		}
	}
}

func TestValidWebhookID(t *testing.T) {
	cases := []struct {
		webhookID      string
		expectedResult bool
	}{
		{
			webhookID:      "e66d248b67c0442fe2cbad7e248651fd4569ee8ecc72ee5a19b0e55ac1ef4492",
			expectedResult: true, // valid, no error
		},
		{
			webhookID:      "E66D248B67C0442FE2CBAD7E248651FD4569EE8ECC72EE5A19B0E55AC1EF4492",
			expectedResult: false, // Invalid hash (capitals A-F instead of a-f)
		},
		{
			webhookID:      "e66d248b67c0442fe2cbad7e248651fd4569ee8ecc72ee5a19b0e55ac1ef449",
			expectedResult: false, // Invalid hash (63 characters)
		},
	}

	for _, c := range cases {
		result := isValidWebhookID(c.webhookID)
		if result != c.expectedResult {
			if c.expectedResult {
				t.Errorf("Expected valid webhook id '%s', but got invalid webhook id", c.webhookID)
			} else {
				t.Errorf("Expected invalid webhook id '%s', but got valid webhook id", c.webhookID)
			}
		}
	}
}
