package server

import (
	"testing"
)

func TestExecuteScript(t *testing.T) {
	webhook := Webhook{}
	payload := []byte("{test: test}")
	err := ExecuteScript(&webhook, payload)
	if err != nil {
		t.Errorf("Expected no error, but got '%+v'", err.Error())
	}
}
