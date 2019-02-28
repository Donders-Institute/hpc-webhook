package server

import (
	"testing"
)

func TestExecuteScript(t *testing.T) {
	webhookID := "550e8400-e29b-41d4-a716-446655440001"
	payload := []byte("{test: test}")
	username := "dccnuser"
	err := ExecuteScript(webhookID, payload, username)
	if err != nil {
		t.Errorf("Expected no error, but got '%+v'", err.Error())
	}
}
