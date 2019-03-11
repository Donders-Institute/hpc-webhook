package server

import (
	"testing"
)

func TestValidConfigurationAddURLPath(t *testing.T) {
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
		result := isValidConfigurationAddURLPath(c.urlPath)
		if result != c.expectedResult {
			if c.expectedResult {
				t.Errorf("Expected valid url path '%s', but got invalid url path", c.urlPath)
			} else {
				t.Errorf("Expected invalid url path '%s', but got valid url path", c.urlPath)
			}
		}
	}
}

func TestValidConfigurationInfoURLPath(t *testing.T) {
	cases := []struct {
		urlPath        string
		expectedResult bool
	}{
		{
			urlPath:        "/configuration/550e8400-e29b-41d4-a716-446655440001",
			expectedResult: true, // Valid configuration URL path, no error
		},
		{
			urlPath:        "/configuration/550e8400-e29b-41d4-a716",
			expectedResult: false, // Invalid hash
		},
		{
			urlPath:        "/configuration/nonexisting",
			expectedResult: false, // Invalid configuration URL path
		},
		{
			urlPath:        "/nonexisting/550e8400-e29b-41d4-a716-446655440001",
			expectedResult: false, // Invalid configuration URL path
		},
	}

	for _, c := range cases {
		result := isValidConfigurationInfoURLPath(c.urlPath)
		if result != c.expectedResult {
			if c.expectedResult {
				t.Errorf("Expected valid url path '%s', but got invalid url path", c.urlPath)
			} else {
				t.Errorf("Expected invalid url path '%s', but got valid url path", c.urlPath)
			}
		}
	}
}

func TestValidConfigurationListURLPath(t *testing.T) {
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
		result := isValidConfigurationListURLPath(c.urlPath)
		if result != c.expectedResult {
			if c.expectedResult {
				t.Errorf("Expected valid url path '%s', but got invalid url path", c.urlPath)
			} else {
				t.Errorf("Expected invalid url path '%s', but got valid url path", c.urlPath)
			}
		}
	}
}

func TestValidConfigurationDeleteURLPath(t *testing.T) {
	cases := []struct {
		urlPath        string
		expectedResult bool
	}{
		{
			urlPath:        "/configuration/550e8400-e29b-41d4-a716-446655440001",
			expectedResult: true, // Valid configuration URL path, no error
		},
		{
			urlPath:        "/configuration/550e8400-e29b-41d4-a716",
			expectedResult: false, // Invalid chas
		},
		{
			urlPath:        "/configuration/nonexisting",
			expectedResult: false, // Invalid configuration URL path
		},
		{
			urlPath:        "/nonexisting/550e8400-e29b-41d4-a716-446655440001",
			expectedResult: false, // Invalid configuration URL path
		},
	}

	for _, c := range cases {
		result := isValidConfigurationDeleteURLPath(c.urlPath)
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
			urlPath:        "/webhook/550e8400-e29b-41d4-a716-446655440001",
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
			urlPath:        "/nonexisting/550e8400-e29b-41d4-a716-446655440001",
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
			webhookID:      "550e8400-e29b-41d4-a716-446655440001",
			expectedResult: true, // valid, no error (i.e. 36 characters with 4 hyphens)
		},
		{
			webhookID:      "550E8400-E29b-41D4-A716-446655440001",
			expectedResult: false, // Invalid hash (i.e. capitals A-F instead of a-f)
		},
		{
			webhookID:      "550e8400-e29b-41d4-a716-44665544000",
			expectedResult: false, // Invalid hash (i.e. 35 characters with 4 hyphens)
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

func TestValidateConfigurationRequest(t *testing.T) {
	cases := []struct {
		conf            ConfigurationRequest
		validateHash    bool
		validateScript  bool
		validateCreated bool
		expectedResult  bool
	}{
		{
			conf: ConfigurationRequest{
				Hash:        "550e8400-e29b-41d4-a716-446655440001",
				Groupname:   "dccngroup",
				Username:    "dccnuser",
				Script:      "script.sh",
				Description: "description",
				Created:     "2019-03-11 10:10:25",
			},
			validateHash:    true,
			validateScript:  true,
			validateCreated: true,
			expectedResult:  true, // valid, no error (i.e. 36 characters with 4 hyphens)
		},
		{
			conf: ConfigurationRequest{
				Hash:        "550E8400-E29b-41D4-A716-446655440001",
				Groupname:   "dccngroup",
				Username:    "dccnuser",
				Script:      "script.sh",
				Description: "description",
				Created:     "2019-03-11 10:10:25",
			},
			validateHash:    true,
			validateScript:  true,
			validateCreated: true,
			expectedResult:  false, // Invalid hash (i.e. capitals A-F instead of a-f)
		},
		{
			conf: ConfigurationRequest{
				Hash:        "550e8400-e29b-41d4-a716-44665544000",
				Groupname:   "dccngroup",
				Username:    "dccnuser",
				Script:      "script.sh",
				Description: "description",
				Created:     "2019-03-11 10:10:25",
			},
			validateHash:    true,
			validateScript:  true,
			validateCreated: true,
			expectedResult:  false, // Invalid hash (i.e. 35 characters with 4 hyphens)
		},
		{
			conf: ConfigurationRequest{
				Hash:        "550e8400-e29b-41d4-a716-446655440001",
				Groupname:   "dccngroup",
				Username:    "",
				Script:      "script.sh",
				Description: "description",
				Created:     "2019-03-11 10:10:25",
			},
			validateHash:    true,
			validateScript:  true,
			validateCreated: true,
			expectedResult:  false, // Invalid username
		},
		{
			conf: ConfigurationRequest{
				Hash:        "550e8400-e29b-41d4-a716-446655440001",
				Groupname:   "",
				Username:    "dccnuser",
				Script:      "script.sh",
				Description: "description",
				Created:     "2019-03-11 10:10:25",
			},
			validateHash:    true,
			validateScript:  true,
			validateCreated: true,
			expectedResult:  false, // Invalid groupname
		},
		{
			conf: ConfigurationRequest{
				Hash:        "",
				Groupname:   "dccngroup",
				Username:    "dccnuser",
				Script:      "script.sh",
				Description: "description",
				Created:     "2019-03-11 10:10:25",
			},
			validateHash:    false,
			validateScript:  true,
			validateCreated: true,
			expectedResult:  true, // Empty hash but no error because validateHash = false
		},
		{
			conf: ConfigurationRequest{
				Hash:        "550e8400-e29b-41d4-a716-446655440001",
				Groupname:   "dccngroup",
				Username:    "dccnuser",
				Script:      "",
				Description: "description",
				Created:     "2019-03-11 10:10:25",
			},
			validateHash:    true,
			validateScript:  false,
			validateCreated: true,
			expectedResult:  true, // Empty script but no error because validateScript = false
		},
		{
			conf: ConfigurationRequest{
				Hash:        "550e8400-e29b-41d4-a716-446655440001",
				Groupname:   "dccngroup",
				Username:    "dccnuser",
				Script:      "script.sh",
				Description: "description",
				Created:     "",
			},
			validateHash:    true,
			validateScript:  true,
			validateCreated: false,
			expectedResult:  true, // Empty creation time but no error because validateCreated = false
		},
	}

	for _, c := range cases {
		err := validateConfigurationRequest(c.conf, c.validateHash, c.validateScript, c.validateCreated)
		if c.expectedResult {
			if err != nil {
				t.Errorf("Expected valid configuration request '%s', but got invalid configuration request", c.conf)
			}
		} else {
			if err == nil {
				t.Errorf("Expected invalid configuration request '%s', but got valid configuration request", c.conf)
			}
		}
	}
}
