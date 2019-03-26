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
	groupname := "dccngroup"
	password := "somepassword"
	dataDir := path.Join("..", "..", "test", "results", "executeScript", "data")
	keyDir := path.Join("..", "..", "test", "results", "executeScript", "keys")
	homeDir := path.Join("..", "..", "test", "results", "executeScript", "home")
	payloadDir := path.Join(dataDir, "payloads", username)
	privateKeyFilename := path.Join(keyDir, "hpc-webhook")
	publicKeyFilename := path.Join(keyDir, "hpc-webhook.pub")
	userScriptDir := path.Join(homeDir, groupname, username, WebhooksWorkDir, webhookID)
	payloadFilename := path.Join(payloadDir, PayLoadName)
	targetPayloadDir := userScriptDir
	targetPayloadFilename := path.Join(targetPayloadDir, PayLoadName)
	userScriptPathFilename := path.Join(userScriptDir, ScriptName)

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

	// Create the key files
	err = os.MkdirAll(keyDir, os.ModePerm)
	if err != nil {
		t.Errorf("Error writing key dir")
	}
	err = ioutil.WriteFile(publicKeyFilename, []byte("test"), 0644)
	if err != nil {
		t.Errorf("Error writing public key")
	}
	err = ioutil.WriteFile(privateKeyFilename, []byte("test"), 0600)
	if err != nil {
		t.Errorf("Error writing private key")
	}
	defer func() {
		err = os.RemoveAll(keyDir) // clean up when done
		if err != nil {
			t.Fatalf("error %s when removing %s dir", err, keyDir)
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

	// Configure the SSH connection
	privateKey, err := ioutil.ReadFile(privateKeyFilename)
	if err != nil {
		t.Errorf("Expected no error, but got '%+v'", err.Error())
	}
	signer, _ := ssh.ParsePrivateKey(privateKey)
	clientConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := fc.NewClient(remote, clientConfig)
	if err != nil {
		t.Errorf("Expected no error, but got '%+v'", err.Error())
	}

	executeConfig := executeConfiguration{
		privateKeyFilename:     privateKeyFilename,
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
		homeDir:                homeDir,
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
	groupname := "dccngroup"
	password := "somepassword"
	dataDir := path.Join("..", "..", "test", "results", "executeScript", "data")
	keyDir := path.Join("..", "..", "test", "results", "executeScript", "keys")
	homeDir := path.Join("..", "..", "test", "results", "executeScript", "home")
	payloadDir := path.Join(dataDir, "payloads", username)
	privateKeyFilename := path.Join(keyDir, "hpc-webhook")
	publicKeyFilename := path.Join(keyDir, "hpc-webhook.pub")
	userScriptDir := path.Join(homeDir, groupname, username, WebhooksWorkDir, webhookID)
	payloadFilename := path.Join(payloadDir, PayLoadName)
	userScriptPathFilename := path.Join(userScriptDir, ScriptName)
	targetPayloadDir := userScriptDir
	targetPayloadFilename := path.Join(targetPayloadDir, PayLoadName)
	payload := []byte("{test: test}")

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

	// Create the key files
	err = os.MkdirAll(keyDir, os.ModePerm)
	if err != nil {
		t.Errorf("Error writing key dir")
	}
	err = ioutil.WriteFile(publicKeyFilename, []byte("test"), 0644)
	if err != nil {
		t.Errorf("Error writing public key")
	}
	err = ioutil.WriteFile(privateKeyFilename, []byte("test"), 0600)
	if err != nil {
		t.Errorf("Error writing private key")
	}
	defer func() {
		err = os.RemoveAll(keyDir) // clean up when done
		if err != nil {
			t.Fatalf("error %s when removing %s dir", err, keyDir)
		}
	}()

	// Create the payload file
	err = os.MkdirAll(payloadDir, os.ModePerm)
	if err != nil {
		t.Errorf("Error writing payload dir")
	}
	err = ioutil.WriteFile(path.Join(payloadDir, PayLoadName), payload, 0644)
	if err != nil {
		t.Errorf("Error writing payload")
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

	executeConfig := executeConfiguration{
		privateKeyFilename:     privateKeyFilename,
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
		homeDir:                homeDir,
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
