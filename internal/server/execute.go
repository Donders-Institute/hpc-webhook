package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

type executeConfiguration struct {
	privateKeyFilename     string
	tempPrivateKeyDir      string
	tempPrivateKeyFilename string
	payloadFilename        string
	targetPayloadDir       string
	targetPayloadFilename  string
	userScriptPathFilename string
	relayNodeName          string
	remoteServer           string
	dataDir                string
	vaultDir               string
	homeDir                string
	webhookID              string
	payload                []byte
	username               string
	groupname              string
	password               string
}

// Copy the src file to dst.
// Any existing file will be overwritten and will not copy file attributes.
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
	command := fmt.Sprintf("cd ~/.qaas/%s/ && qsub -F %s %s", conf.webhookID, conf.targetPayloadFilename, userScriptFilename)
	fmt.Println(command)
	err = c.Run(session, command)
	if err != nil {
		return err
	}
	return err
}

// ExecuteScript triggers a qsub command on the HPC cluster
func ExecuteScript(c Connector, conf executeConfiguration) error {
	// Copy the private key to the vault directory and
	// change its file permissions to root read and write access only
	err := os.MkdirAll(conf.tempPrivateKeyDir, os.ModePerm)
	if err != nil {
		return err
	}
	fmt.Println(conf.privateKeyFilename)
	fmt.Println(conf.tempPrivateKeyFilename)
	err = copyFile(conf.privateKeyFilename, conf.tempPrivateKeyFilename)
	if err != nil {
		return err
	}
	err = os.Chmod(conf.tempPrivateKeyFilename, 0600)
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
		User: conf.username,
		// Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		Auth: []ssh.AuthMethod{ssh.Password(conf.password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// Start an SSH session on the relay node
	remoteServer := fmt.Sprintf("%s:22", conf.relayNodeName)
	client, err := c.NewClient(remoteServer, clientConfig)
	if err != nil {
		return err
	}

	// Copy the payload to QaaS folder
	err = os.MkdirAll(conf.targetPayloadDir, os.ModePerm)
	if err != nil {
		return err
	}
	err = copyFile(conf.payloadFilename, conf.targetPayloadFilename)
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
