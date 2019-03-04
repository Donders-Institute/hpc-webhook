package server

import (
	"fmt"
	"net"

	"golang.org/x/crypto/ssh"
)

// Connector is an interface to be able to mock SSH connections
type Connector interface {
	NewClient(remote string, clientConfig *ssh.ClientConfig) (*ssh.Client, error)
	NewSession(client *ssh.Client) (*ssh.Session, error)
	Run(session *ssh.Session, command string) error
	CombinedOutput(session *ssh.Session, command string) ([]byte, error)
}

// SSHConnector is used tp replace the standard SSH library functions
type SSHConnector struct {
	description string
}

// NewClient makes it possible to mock SSH dial
func (sc SSHConnector) NewClient(remote string, clientConfig *ssh.ClientConfig) (*ssh.Client, error) {
	return ssh.Dial("tcp", remote, clientConfig)
}

// NewSession makes it possible to mock a SSH session
func (sc SSHConnector) NewSession(client *ssh.Client) (*ssh.Session, error) {
	return client.NewSession()
}

// Run makes it possible to mock a remote Run command
func (sc SSHConnector) Run(session *ssh.Session, command string) error {
	return session.Run(command)
}

// CombinedOutput makes it possible to mock a local CombinedOutput
func (sc SSHConnector) CombinedOutput(session *ssh.Session, command string) ([]byte, error) {
	return session.CombinedOutput(command)
}

func connect(c Connector) error {
	var session *ssh.Session
	var out []byte

	remote := "relaynode.dccn.nl:22"
	username := "dccnuser"
	password := "somepassword"
	clientConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err := c.NewClient(remote, clientConfig)
	if err != nil {
		return err
	}

	session, err = c.NewSession(client)
	if err != nil {
		return err
	}
	defer session.Close()

	command := "ls -1"
	err = c.Run(session, command)

	command = `ssh dccnuser@mentat001.dccn.nl "ls -1"`
	out, err = c.CombinedOutput(session, command)
	if err != nil {
		return err
	}
	fmt.Println(string(out))

	return err
}
