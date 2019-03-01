package server

import (
	"fmt"
	"path"
)

// ExecuteScript triggers a qsub command on the HPC cluster
func ExecuteScript(relayNodeName string, dataDir string, webhookID string, payload []byte, username string) error {
	privateKeyFilename := path.Join(dataDir, "keys", username, "id_rsa")
	payloadFilename := path.Join(dataDir, "payloads", username, "payload")

	fmt.Printf("Execute: %+v\n", webhookID)
	fmt.Printf("Execute: %+v\n", string(payload))
	fmt.Printf("Execute: %+v\n", username)
	fmt.Printf("Private key filename: %+v\n", privateKeyFilename)

	// First copy the payload to QaaS folder
	scpCommand := fmt.Sprintf(`scp -i %s %s %s@%s:~/.qaas/%s/payload`,
		privateKeyFilename, payloadFilename, username, relayNodeName, webhookID)
	fmt.Printf("%s\n", scpCommand)

	// Next, trigger the qsub command
	command := fmt.Sprintf("echo ~/.qaas/%s/script.sh payload | qsub", webhookID)
	sshCommand := fmt.Sprintf(`ssh -i %s %s@%s "%s"`, privateKeyFilename, username, relayNodeName, command)
	fmt.Printf("%s\n", sshCommand)

	return nil
}
