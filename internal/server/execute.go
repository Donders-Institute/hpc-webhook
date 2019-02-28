package server

import (
	"errors"
	"fmt"
)

// TODO: ExecuteScript triggers a qsub command on the HPC cluster
func ExecuteScript(webhookID string, payload []byte, username string) error {
	// var err error
	relayNodeName := "relay-node.dccn.nl"
	privateKeyFilename := "./data/keys/username/id_rsa"
	// publicKeyFilename := "./data/keys/username/id_rsa.pub"

	fmt.Printf("Execute: %+v\n", webhookID)
	fmt.Printf("Execute: %+v\n", string(payload))
	fmt.Printf("Execute: %+v\n", username)
	fmt.Printf("Private key filename: %+v\n", privateKeyFilename)

	command := fmt.Sprintf("echo ~/.qaas/%s/script.sh data | qsub", webhookID)
	sshCommand := fmt.Sprintf(`ssh -i %s %s@%s "%s"`, privateKeyFilename, username, relayNodeName, command)

	fmt.Printf("%s\n", sshCommand)
	return errors.New("execute failure")
	// return err
}
