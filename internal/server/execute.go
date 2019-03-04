package server

import (
	"fmt"
	"io"
	"net"
	"os"
	"path"

	"golang.org/x/crypto/ssh"
)

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func Copy(src string, dst string) error {
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

func copyPayload(privateKeyFilename string, payloadFilename string, username string, password string, relayNodeName string, webhookID string) error {

	// Copy the private key to the vault
	tempDir := path.Join("/vault", username)
	tempPrivateKeyFilename := path.Join(tempDir, "id_rsa")
	err := os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		return err
	}
	err = Copy(privateKeyFilename, tempPrivateKeyFilename)
	if err != nil {
		return err
	}
	err = os.Chmod(tempPrivateKeyFilename, 0600)
	if err != nil {
		return err
	}

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

	scpCommand := fmt.Sprintf(`scp -i %s %s %s@%s:~/.qaas/%s/payload`,
		tempPrivateKeyFilename, payloadFilename, username, relayNodeName, webhookID)
	fmt.Println(scpCommand)
	err = session.Run(scpCommand)
	if err != nil {
		fmt.Println("Warning: something went wrong copying the payload. Skipping ...")
		return nil
	}

	return err
}

func triggerQsubCommand(privateKeyFilename string, payloadFilename string, username string, password string, relayNodeName string, webhookID string) error {
	// privateKey, err := ioutil.ReadFile(privateKeyFilename)
	// if err != nil {
	// 	return err
	// }

	//signer, _ := ssh.ParsePrivateKey(privateKey)
	clientConfig := &ssh.ClientConfig{
		User: username,
		//Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		Auth: []ssh.AuthMethod{ssh.Password(password)},
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

	command := fmt.Sprintf("cd ~/.qaas/%s/ && echo ~/.qaas/%s/script.sh payload | qsub", webhookID, webhookID)
	// sshCommand := fmt.Sprintf(`ssh -i %s %s@%s "%s"`, privateKeyFilename, username, relayNodeName, command)
	var out []byte
	out, err = session.CombinedOutput(command)
	if err != nil {
		return err
	}
	fmt.Println(string(out))

	return err
}

// ExecuteScript triggers a qsub command on the HPC cluster
func ExecuteScript(relayNodeName string, dataDir string, webhookID string, payload []byte, username string, password string) error {
	privateKeyFilename := path.Join(dataDir, "keys", username, "id_rsa")
	payloadFilename := path.Join(dataDir, "payloads", username, "payload")

	// First copy the payload to QaaS folder
	err := copyPayload(privateKeyFilename, payloadFilename, username, password, relayNodeName, webhookID)
	if err != nil {
		return err
	}

	// Next, trigger the qsub command
	err = triggerQsubCommand(privateKeyFilename, payloadFilename, username, password, relayNodeName, webhookID)
	if err != nil {
		return err
	}

	return err
}
