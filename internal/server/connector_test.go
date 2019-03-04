package server

import (
	"testing"

	"golang.org/x/crypto/ssh"
)

type FakeConnector struct {
	description string
}

func (f FakeConnector) NewClient(remote string, clientConfig *ssh.ClientConfig) (*ssh.Client, error) {
	return nil, nil
}

func (f FakeConnector) NewSession(client *ssh.Client) (*ssh.Session, error) {
	return nil, nil
}

func (f FakeConnector) Run(session *ssh.Session, command string) error {
	return nil
}

func (f FakeConnector) CombinedOutput(session *ssh.Session, command string) ([]byte, error) {
	return nil, nil
}

func TestConnect(t *testing.T) {
	fc := FakeConnector{
		description: "fake SSH connection",
	}
	err := connect(fc)
	if err != nil {
		t.Errorf("Expected no error, but got '%+v'", err.Error())
	}
}
