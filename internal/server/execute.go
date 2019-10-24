package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

type executeConfiguration struct {
	privateKeyFilename       string
	payloadFilename          string
	targetPayloadDir         string
	targetPayloadFilename    string
	userScriptPathFilename   string
	relayNodeName            string
	connectionTimeoutSeconds int
	remoteServer             string
	dataDir                  string
	homeDir                  string
	webhookID                string
	payload                  []byte
	username                 string
	groupname                string
	password                 string
}

// CopyFile copies a source file to a destination file.
// Any existing file will be overwritten and will not copy file attributes.
func CopyFile(src string, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func triggerQsubCommand(c Connector, client *ssh.Client, conf executeConfiguration) error {
	session, err := c.NewSession(client)
	if err != nil {
		return err
	}
	defer c.CloseSession(session)

	// Grab the path to the user script
	contents, err := ioutil.ReadFile(conf.userScriptPathFilename)
	if err != nil {
		return err
	}
	userScriptFilename := string(contents)

	// Go the correct folder and run the qsub command from there
	command := fmt.Sprintf(`bash -l -c "cd ~/%s/%s/ && qsub -F %s %s"`, WebhooksWorkDir, conf.webhookID, conf.targetPayloadFilename, userScriptFilename)
	fmt.Println(command)
	err = c.Run(session, command)
	if err != nil {
		return err
	}
	return err
}

// ExecuteScript triggers a qsub command on the HPC cluster
func ExecuteScript(c Connector, conf executeConfiguration) error {
	// Configure the SSH connection
	privateKey, err := ioutil.ReadFile(conf.privateKeyFilename)
	if err != nil {
		return err
	}
	signer, _ := ssh.ParsePrivateKey(privateKey)
	clientConfig := &ssh.ClientConfig{
		User: conf.username,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: time.Duration(conf.connectionTimeoutSeconds) * time.Second,
	}

	// Start an SSH session on the relay node
	remoteServer := fmt.Sprintf("%s:22", conf.relayNodeName)
	client, err := c.NewClient(remoteServer, clientConfig)
	if err != nil {
		return err
	}
	defer c.CloseConnection(client)

	// Copy the payload to HPC webhooks folder
	err = os.MkdirAll(conf.targetPayloadDir, os.ModePerm)
	if err != nil {
		return err
	}
	err = CopyFile(conf.payloadFilename, conf.targetPayloadFilename)
	if err != nil {
		return err
	}

	// Trigger the qsub command
	err = triggerQsubCommand(c, client, conf)
	if err != nil {
		return err
	}

	return err
}
