package server

import (
	"io/ioutil"
	"net"
	"os"
	"path"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestCopyPayload(t *testing.T) {
	var session *ssh.Session

	relayNodeName := "relaynode.dccn.nl"
	remote := net.JoinHostPort(relayNodeName, "22")
	webhookID := "550e8400-e29b-41d4-a716-446655440001"
	username := "dccnuser"
	password := "somepassword"
	testDir := path.Join("..", "..", "test", "results", "executeScript")
	keyDir := path.Join(testDir, "keys", username)
	payloadDir := path.Join(testDir, "payloads", username)
	tempPrivateKeyFilename := path.Join(keyDir, "id_rsa")
	payloadFilename := path.Join(payloadDir, "payload")

	executeConfig := executeConfiguration{
		tempPrivateKeyFilename: tempPrivateKeyFilename,
		payloadFilename:        payloadFilename,
		username:               username,
		password:               password,
		relayNodeName:          relayNodeName,
		webhookID:              webhookID,
	}

	clientConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	fc := FakeConnector{
		Description: "fake SSH connection",
	}

	client, err := fc.NewClient(remote, clientConfig)
	if err != nil {
		t.Errorf("Expected no error, but got '%+v'", err.Error())
	}

	session, err = fc.NewSession(client)
	if err != nil {
		t.Errorf("Expected no error, but got '%+v'", err.Error())
	}
	defer fc.CloseSession(session)

	err = copyPayload(fc, session, executeConfig)
	if err != nil {
		t.Errorf("Expected no error, but got '%+v'", err.Error())
	}
}

func TestTriggerQsubCommand(t *testing.T) {
	var session *ssh.Session

	relayNodeName := "relaynode.dccn.nl"
	remote := net.JoinHostPort(relayNodeName, "22")
	webhookID := "550e8400-e29b-41d4-a716-446655440001"
	username := "dccnuser"
	password := "somepassword"
	testDir := path.Join("..", "..", "test", "results", "executeScript")
	keyDir := path.Join(testDir, "keys", username)
	payloadDir := path.Join(testDir, "payloads", username)
	tempPrivateKeyFilename := path.Join(keyDir, "id_rsa")
	payloadFilename := path.Join(payloadDir, "payload")

	executeConfig := executeConfiguration{
		tempPrivateKeyFilename: tempPrivateKeyFilename,
		payloadFilename:        payloadFilename,
		username:               username,
		password:               password,
		relayNodeName:          relayNodeName,
		webhookID:              webhookID,
	}

	clientConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	fc := FakeConnector{
		Description: "fake SSH connection",
	}

	client, err := fc.NewClient(remote, clientConfig)
	if err != nil {
		t.Errorf("Expected no error, but got '%+v'", err.Error())
	}

	session, err = fc.NewSession(client)
	if err != nil {
		t.Errorf("Expected no error, but got '%+v'", err.Error())
	}
	defer fc.CloseSession(session)

	err = triggerQsubCommand(fc, session, executeConfig)
	if err != nil {
		t.Errorf("Expected no error, but got '%+v'", err.Error())
	}
}

func TestExecuteScript(t *testing.T) {
	relayNode := "relaynode.dccn.nl"
	webhookID := "550e8400-e29b-41d4-a716-446655440001"
	payload := []byte("{test: test}")
	username := "dccnuser"
	password := "somepassword"
	testDir := path.Join("..", "..", "test", "results", "executeScript")
	keyDir := path.Join(testDir, "keys", username)

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

	// Create the payload file
	err = ioutil.WriteFile(path.Join(testDir, "payload"), payload, 0644)
	if err != nil {
		t.Errorf("Error writing test result")
	}

	// Execute the script
	fc := FakeConnector{
		Description: "fake SSH connection",
	}
	err = ExecuteScript(fc, relayNode, testDir, webhookID, payload, username, password)
	if err != nil {
		t.Errorf("Expected no error, but got '%+v'", err.Error())
	}
}
