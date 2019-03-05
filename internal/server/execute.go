package server

import (
	"fmt"
	"io"
	"net"
	"os"
	"path"

	"golang.org/x/crypto/ssh"
)

type executeConfiguration struct {
	tempPrivateKeyFilename string
	payloadFilename        string
	username               string
	password               string
	relayNodeName          string
	webhookID              string
}

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func copyFile(src string, dst string) error {
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

func copyPayload(c Connector, client *ssh.Client, conf executeConfiguration) error {
	session, err := c.NewSession(client)
	if err != nil {
		return err
	}
	defer c.CloseSession(session)

	scpCommand := fmt.Sprintf(`scp -i %s %s %s@%s:~/.qaas/%s/payload`,
		conf.tempPrivateKeyFilename, conf.payloadFilename, conf.username, conf.relayNodeName, conf.webhookID)
	fmt.Println(scpCommand)
	out, err := c.CombinedOutput(session, scpCommand)
	if err != nil {
		fmt.Println("Warning: something went wrong copying the payload. Skipping ...")
		return nil
	}
	fmt.Println(string(out))
	return err
}

func triggerQsubCommand(c Connector, client *ssh.Client, conf executeConfiguration) error {
	session, err := c.NewSession(client)
	if err != nil {
		return err
	}
	defer c.CloseSession(session)

	command := fmt.Sprintf("cd ~/.qaas/%s/ && cat ~/.qaas/%s/script.sh payload | qsub", conf.webhookID, conf.webhookID)
	err = c.Run(session, command)
	if err != nil {
		return err
	}
	return err
}

// ExecuteScript triggers a qsub command on the HPC cluster
func ExecuteScript(c Connector, relayNodeName string, dataDir string, webhookID string, payload []byte, username string, password string) error {
	privateKeyFilename := path.Join(dataDir, "keys", username, "id_rsa")
	payloadFilename := path.Join(dataDir, "payloads", username, "payload")
	remote := fmt.Sprintf("%s:22", relayNodeName)

	// Copy the private key to the vault directroy and
	// change its file permissions to root read and write access only
	tempDir := path.Join("/vault", username)
	tempPrivateKeyFilename := path.Join(tempDir, "id_rsa")
	err := os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		return err
	}
	err = copyFile(privateKeyFilename, tempPrivateKeyFilename)
	if err != nil {
		return err
	}
	err = os.Chmod(tempPrivateKeyFilename, 0600)
	if err != nil {
		return err
	}

	// Configure the SSH connection
	// privateKey, err := ioutil.ReadFile(privateKeyFilename)
	// if err != nil {
	// 	return err
	// }
	//signer, _ := ssh.ParsePrivateKey(privateKey)
	clientConfig := &ssh.ClientConfig{
		User: username,
		// Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// Start an SSH session on the relay node
	client, err := c.NewClient(remote, clientConfig)
	if err != nil {
		return err
	}

	// Combine the parameters
	conf := executeConfiguration{
		tempPrivateKeyFilename: tempPrivateKeyFilename,
		payloadFilename:        payloadFilename,
		username:               username,
		password:               password,
		relayNodeName:          relayNodeName,
		webhookID:              webhookID,
	}

	// Copy the payload to QaaS folder
	err = copyPayload(c, client, conf)
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
