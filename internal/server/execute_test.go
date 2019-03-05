package server

import (
	"io/ioutil"
	"net"
	"os"
	"path"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestTriggerQsubCommand(t *testing.T) {
	relayNodeName := "relaynode.dccn.nl"
	remote := net.JoinHostPort(relayNodeName, "22")
	webhookID := "550e8400-e29b-41d4-a716-446655440001"
	username := "dccnuser"
	groupname := "tg"
	password := "somepassword"
	dataDir := path.Join("..", "..", "test", "results", "executeScript", "data")
	vaultDir := path.Join("..", "..", "test", "results", "executeScript", "vault")
	homeDir := path.Join("..", "..", "test", "results", "executeScript", "home")
	payloadDir := path.Join(dataDir, "payloads", username)
	privateKeyFilename := path.Join(dataDir, "keys", username, "id_rsa")
	userScriptDir := path.Join(homeDir, groupname, username, ".qaas", webhookID)
	tempPrivateKeyDir := path.Join(vaultDir, username, "id_rsa")
	tempPrivateKeyFilename := path.Join(tempPrivateKeyDir, "id_rsa")
	payloadFilename := path.Join(payloadDir, "payload")
	targetPayloadDir := userScriptDir
	targetPayloadFilename := path.Join(targetPayloadDir, "payload")
	userScriptPathFilename := path.Join(userScriptDir, "script.sh")

	executeConfig := executeConfiguration{
		privateKeyFilename:     privateKeyFilename,
		tempPrivateKeyDir:      tempPrivateKeyDir,
		tempPrivateKeyFilename: tempPrivateKeyFilename,
		payloadFilename:        payloadFilename,
		targetPayloadDir:       targetPayloadDir,
		targetPayloadFilename:  targetPayloadFilename,
		userScriptPathFilename: userScriptPathFilename,
		username:               username,
		groupname:              groupname,
		password:               password,
		relayNodeName:          relayNodeName,
		webhookID:              webhookID,
		dataDir:                dataDir,
		vaultDir:               vaultDir,
		homeDir:                homeDir,
	}

	clientConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// Create the data dir
	err := os.MkdirAll(dataDir, os.ModePerm)
	if err != nil {
		t.Errorf("Error writing test result data dir")
	}
	defer func() {
		err = os.RemoveAll(dataDir) // clean up when done
		if err != nil {
			t.Fatalf("error %s when removing %s dir", err, dataDir)
		}
	}()

	// Create the home dir
	err = os.MkdirAll(homeDir, os.ModePerm)
	if err != nil {
		t.Errorf("Error writing test result home dir")
	}
	defer func() {
		err = os.RemoveAll(homeDir) // clean up when done
		if err != nil {
			t.Fatalf("error %s when removing %s dir", err, homeDir)
		}
	}()

	// Create the user script file
	err = os.MkdirAll(userScriptDir, os.ModePerm)
	if err != nil {
		t.Errorf("Error writing user script dir")
	}
	err = ioutil.WriteFile(userScriptPathFilename, []byte("test.sh"), 0644)
	if err != nil {
		t.Errorf("Error writing script.sh")
	}

	fc := FakeConnector{
		Description: "fake SSH connection",
	}

	client, err := fc.NewClient(remote, clientConfig)
	if err != nil {
		t.Errorf("Expected no error, but got '%+v'", err.Error())
	}

	err = triggerQsubCommand(fc, client, executeConfig)
	if err != nil {
		t.Errorf("Expected no error, but got '%+v'", err.Error())
	}
}

func TestExecuteScript(t *testing.T) {
	relayNodeName := "relaynode.dccn.nl"
	remoteServer := net.JoinHostPort(relayNodeName, "22")
	webhookID := "550e8400-e29b-41d4-a716-446655440001"
	username := "dccnuser"
	groupname := "tg"
	password := "somepassword"
	dataDir := path.Join("..", "..", "test", "results", "executeScript", "data")
	vaultDir := path.Join("..", "..", "test", "results", "executeScript", "vault")
	homeDir := path.Join("..", "..", "test", "results", "executeScript", "home")
	keyDir := path.Join(dataDir, "keys", username)
	payloadDir := path.Join(dataDir, "payloads", username)
	privateKeyFilename := path.Join(dataDir, "keys", username, "id_rsa")
	userScriptDir := path.Join(homeDir, groupname, username, ".qaas", webhookID)
	tempPrivateKeyDir := path.Join(vaultDir, username, "id_rsa")
	tempPrivateKeyFilename := path.Join(tempPrivateKeyDir, "id_rsa")
	payloadFilename := path.Join(payloadDir, "payload")
	userScriptPathFilename := path.Join(userScriptDir, "script.sh")
	targetPayloadDir := userScriptDir
	targetPayloadFilename := path.Join(targetPayloadDir, "payload")
	payload := []byte("{test: test}")

	executeConfig := executeConfiguration{
		privateKeyFilename:     privateKeyFilename,
		tempPrivateKeyDir:      tempPrivateKeyDir,
		tempPrivateKeyFilename: tempPrivateKeyFilename,
		payloadFilename:        payloadFilename,
		targetPayloadDir:       targetPayloadDir,
		targetPayloadFilename:  targetPayloadFilename,
		userScriptPathFilename: userScriptPathFilename,
		username:               username,
		groupname:              groupname,
		password:               password,
		relayNodeName:          relayNodeName,
		remoteServer:           remoteServer,
		webhookID:              webhookID,
		dataDir:                dataDir,
		vaultDir:               vaultDir,
		homeDir:                homeDir,
	}

	// Create the data dir
	err := os.MkdirAll(dataDir, os.ModePerm)
	if err != nil {
		t.Errorf("Error writing test result data dir")
	}
	// defer func() {
	// 	err = os.RemoveAll(dataDir) // clean up when done
	// 	if err != nil {
	// 		t.Fatalf("error %s when removing %s dir", err, dataDir)
	// 	}
	// }()

	// Create the home dir
	err = os.MkdirAll(homeDir, os.ModePerm)
	if err != nil {
		t.Errorf("Error writing test result home dir")
	}
	// defer func() {
	// 	err = os.RemoveAll(homeDir) // clean up when done
	// 	if err != nil {
	// 		t.Fatalf("error %s when removing %s dir", err, homeDir)
	// 	}
	// }()

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
	err = os.MkdirAll(payloadDir, os.ModePerm)
	if err != nil {
		t.Errorf("Error writing test result dir")
	}
	err = ioutil.WriteFile(path.Join(payloadDir, "payload"), payload, 0644)
	if err != nil {
		t.Errorf("Error writing test result")
	}

	// Create the user script file
	err = os.MkdirAll(userScriptDir, os.ModePerm)
	if err != nil {
		t.Errorf("Error writing user script dir")
	}
	err = ioutil.WriteFile(userScriptPathFilename, []byte("test.sh"), 0644)
	if err != nil {
		t.Errorf("Error writing script.sh")
	}

	// Execute the script
	fc := FakeConnector{
		Description: "fake SSH connection",
	}
	err = ExecuteScript(fc, executeConfig)
	if err != nil {
		t.Errorf("Expected no error, but got '%+v'", err.Error())
	}
}
