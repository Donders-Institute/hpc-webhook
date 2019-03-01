package server

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestExecuteScript(t *testing.T) {
	relayNode := "relaynode.dccn.nl"
	testDir := path.Join("..", "..", "test", "results", "executeScript")
	webhookID := "550e8400-e29b-41d4-a716-446655440001"
	payload := []byte("{test: test}")
	username := "dccnuser"

	// Create the test results dir if it does not exist
	err := os.MkdirAll(testDir, os.ModePerm)
	if err != nil {
		t.Errorf("Error writing test result dir")
	}
	defer func() {
		err = os.RemoveAll(testDir) // cleanup when done
		if err != nil {
			t.Fatalf("error %s when removing %s dir", err, testDir)
		}
	}()

	err = ioutil.WriteFile(path.Join(testDir, "payload"), payload, 0644)
	if err != nil {
		t.Errorf("Error writing test result")
	}

	err = ExecuteScript(relayNode, testDir, webhookID, payload, username)
	if err != nil {
		t.Errorf("Expected no error, but got '%+v'", err.Error())
	}
}
