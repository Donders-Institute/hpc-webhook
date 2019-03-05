package server

import (
	"fmt"
	"net"
	"testing"

	"golang.org/x/crypto/ssh"
)

type FakeConnector struct {
	description string
}

func (fc FakeConnector) NewClient(remote string, clientConfig *ssh.ClientConfig) (*ssh.Client, error) {
	return nil, nil
}

func (fc FakeConnector) NewSession(client *ssh.Client) (*ssh.Session, error) {
	return nil, nil
}

func (fc FakeConnector) Run(session *ssh.Session, command string) error {
	return nil
}

func (fc FakeConnector) CombinedOutput(session *ssh.Session, command string) ([]byte, error) {
	return nil, nil
}

func (fc FakeConnector) CloseSession(session *ssh.Session) error {
	var err error
	return err
}

func TestConnect(t *testing.T) {
	fc := FakeConnector{
		description: "fake SSH connection",
	}
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
	client, err := fc.NewClient(remote, clientConfig)
	if err != nil {
		t.Errorf("Expected no error, but got '%+v'", err.Error())
	}

	session, err = fc.NewSession(client)
	if err != nil {
		t.Errorf("Expected no error, but got '%+v'", err.Error())
	}
	defer fc.CloseSession(session)

	command := "ls -1"
	err = fc.Run(session, command)

	command = `ssh dccnuser@relaynode.dccn.nl "ls -1"`
	out, err = fc.CombinedOutput(session, command)
	if err != nil {
		t.Errorf("Expected no error, but got '%+v'", err.Error())
	}
	fmt.Println(string(out))
}
