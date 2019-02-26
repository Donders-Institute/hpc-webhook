package server

import (
	"fmt"
)

// ExecuteScript triggers a qsub command on the Torque cluster
func ExecuteScript(webhook *Webhook) error {
	var err error
	fmt.Printf("Execute: '%+v'\n", webhook)
	return err
}
