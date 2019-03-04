package server

import (
	"fmt"
	"io/ioutil"
	"net"
	"path"

	"golang.org/x/crypto/ssh"
)

func copyPayload(privateKeyFilename string, payloadFilename string, username string, relayNodeName string, webhookID string) error {
	privateKey, err := ioutil.ReadFile(privateKeyFilename)
	if err != nil {
		return err
	}

	signer, _ := ssh.ParsePrivateKey(privateKey)
	clientConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	target := fmt.Sprintf("%s:22", relayNodeName)
	client, err := ssh.Dial("tcp", target, clientConfig)
	if err != nil {
		return err
	}
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	scpCommand := fmt.Sprintf(`scp %s %s@%s:~/.qaas/%s/payload`,
		payloadFilename, username, relayNodeName, webhookID)
	err = session.Run(scpCommand)
	if err != nil {
		return err
	}

	return err
}

func triggerQsubCommand(privateKeyFilename string, payloadFilename string, username string, relayNodeName string, webhookID string) error {
	privateKey, err := ioutil.ReadFile(privateKeyFilename)
	if err != nil {
		return err
	}

	signer, _ := ssh.ParsePrivateKey(privateKey)
	clientConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	target := fmt.Sprintf("%s:22", relayNodeName)
	client, err := ssh.Dial("tcp", target, clientConfig)
	if err != nil {
		return err
	}
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	scpCommand := fmt.Sprintf(`scp %s %s@%s:~/.qaas/%s/payload`,
		payloadFilename, username, relayNodeName, webhookID)
	err = session.Run(scpCommand)
	if err != nil {
		return err
	}

	command := fmt.Sprintf("echo ~/.qaas/%s/script.sh payload | qsub", webhookID)
	sshCommand := fmt.Sprintf(`ssh -i %s %s@%s "%s"`, privateKeyFilename, username, relayNodeName, command)
	err = session.Run(sshCommand)
	if err != nil {
		return err
	}

	return err
}

// ExecuteScript triggers a qsub command on the HPC cluster
func ExecuteScript(relayNodeName string, dataDir string, webhookID string, payload []byte, username string) error {
	privateKeyFilename := path.Join(dataDir, "keys", username, "id_rsa")
	payloadFilename := path.Join(dataDir, "payloads", username, "payload")

	// First copy the payload to QaaS folder
	err := copyPayload(privateKeyFilename, payloadFilename, username, relayNodeName, webhookID)
	if err != nil {
		return err
	}

	// Next, trigger the qsub command
	err = triggerQsubCommand(privateKeyFilename, payloadFilename, username, relayNodeName, webhookID)
	if err != nil {
		return err
	}

	return err
}
