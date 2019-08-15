package server

import (
	"golang.org/x/crypto/ssh"
)

// Connector is an interface to be able to mock SSH connections
type Connector interface {
	NewClient(remoteServer string, clientConfig *ssh.ClientConfig) (*ssh.Client, error)
	NewSession(client *ssh.Client) (*ssh.Session, error)
	Run(session *ssh.Session, command string) error
	CombinedOutput(session *ssh.Session, command string) ([]byte, error)
	CloseSession(session *ssh.Session) error
	CloseConnection(client *ssh.Client)
}

// SSHConnector is used tp replace the standard SSH library functions
type SSHConnector struct {
	Description string
}

// NewClient makes it possible to mock SSH dial
func (c SSHConnector) NewClient(remoteServer string, clientConfig *ssh.ClientConfig) (*ssh.Client, error) {
	return ssh.Dial("tcp", remoteServer, clientConfig)
}

// NewSession makes it possible to mock a SSH session
func (c SSHConnector) NewSession(client *ssh.Client) (*ssh.Session, error) {
	return client.NewSession()
}

// Run makes it possible to mock a remote Run command
func (c SSHConnector) Run(session *ssh.Session, command string) error {
	return session.Run(command)
}

// CombinedOutput makes it possible to mock a local CombinedOutput
func (c SSHConnector) CombinedOutput(session *ssh.Session, command string) ([]byte, error) {
	return session.CombinedOutput(command)
}

// CloseSession makes it possible to mock the closing of a session
func (c SSHConnector) CloseSession(session *ssh.Session) error {
	return session.Close()
}

// CloseConnection makes it possible to mock the closing of a session
func (c SSHConnector) CloseConnection(client *ssh.Client) {
	client.Conn.Close()
}
