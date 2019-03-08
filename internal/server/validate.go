package server

import (
	"errors"
	"fmt"
	"regexp"
)

var validConfigurationAddURLPathRegexString = fmt.Sprintf(`^%s$`, ConfigurationPath)
var validConfigurationAddURLPathRegex = regexp.MustCompile(validConfigurationAddURLPathRegexString)

var validConfigurationInfoURLPathRegexString = fmt.Sprintf(`^%s/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, ConfigurationPath)
var validConfigurationInfoURLPathRegex = regexp.MustCompile(validConfigurationInfoURLPathRegexString)

var validConfigurationListURLPathRegexString = fmt.Sprintf(`^%s$`, ConfigurationPath)
var validConfigurationListURLPathRegex = regexp.MustCompile(validConfigurationListURLPathRegexString)

var validConfigurationDeleteURLPathRegexString = fmt.Sprintf(`^%s/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, ConfigurationPath)
var validConfigurationDeleteURLPathRegex = regexp.MustCompile(validConfigurationDeleteURLPathRegexString)

var validURLPathRegexString = fmt.Sprintf(`^%s[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, WebhookPath)
var validURLPathRegex = regexp.MustCompile(validURLPathRegexString)

var validWebhookIDRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func isValidConfigurationAddURLPath(urlPath string) bool {
	return validConfigurationAddURLPathRegex.MatchString(urlPath)
}

func isValidConfigurationInfoURLPath(urlPath string) bool {
	return validConfigurationInfoURLPathRegex.MatchString(urlPath)
}

func isValidConfigurationListURLPath(urlPath string) bool {
	return validConfigurationListURLPathRegex.MatchString(urlPath)
}

func isValidConfigurationDeleteURLPath(urlPath string) bool {
	return validConfigurationDeleteURLPathRegex.MatchString(urlPath)
}

func isValidURLPath(urlPath string) bool {
	return validURLPathRegex.MatchString(urlPath)
}

func isValidWebhookID(webhookID string) bool {
	return validWebhookIDRegex.MatchString(webhookID)
}

func validateConfigurationRequest(conf ConfigurationRequest, validateHash bool) error {
	if validateHash && !isValidWebhookID(conf.Hash) {
		return errors.New("invalid configuration request: invalid hash")
	}
	if conf.Username == "" {
		return errors.New("invalid configuration request: username missing")
	}
	if conf.Groupname == "" {
		return errors.New("invalid configuration request: groupname missing")
	}
	return nil
}
