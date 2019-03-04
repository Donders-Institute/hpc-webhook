package server

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestExecuteScript(t *testing.T) {
	relayNode := "relaynode.dccn.nl"
	webhookID := "550e8400-e29b-41d4-a716-446655440001"
	payload := []byte("{test: test}")
	username := "dccnuser"
	password := "somepassword"
	testDir := path.Join("..", "..", "test", "results", "executeScript")
	keyDir := path.Join("..", "..", "test", "results", "executeScript", "keys", username)

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

	// Create the key files
	err = os.MkdirAll(keyDir, os.ModePerm)
	if err != nil {
		t.Errorf("Error writing test result dir")
	}
	err = ioutil.WriteFile(path.Join(keyDir, "id_rsa.pub"), []byte("test"), 0644)
	if err != nil {
		t.Errorf("Error writing test result")
	}
	err = ioutil.WriteFile(path.Join(keyDir, "id_rsa"), []byte("test"), 0600)
	if err != nil {
		t.Errorf("Error writing test result")
	}

	err = ioutil.WriteFile(path.Join(testDir, "payload"), payload, 0644)
	if err != nil {
		t.Errorf("Error writing test result")
	}

	err = ExecuteScript(relayNode, testDir, webhookID, payload, username, password)
	if err != nil {
		t.Errorf("Expected no error, but got '%+v'", err.Error())
	}
}
