package server

import (
	"fmt"
	"regexp"
)

var validConfigurationURLPathRegexString = fmt.Sprintf(`^%s$`, ConfigurationPath)
var validConfigurationURLPathRegex = regexp.MustCompile(validConfigurationURLPathRegexString)

var validURLPathRegexString = fmt.Sprintf(`^%s[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, WebhookPath)
var validURLPathRegex = regexp.MustCompile(validURLPathRegexString)

var validWebhookIDRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func isValidConfigurationURLPath(urlPath string) bool {
	return validConfigurationURLPathRegex.MatchString(urlPath)
}

func isValidURLPath(urlPath string) bool {
	return validURLPathRegex.MatchString(urlPath)
}

func isValidWebhookID(webhookID string) bool {
	return validWebhookIDRegex.MatchString(webhookID)
}
