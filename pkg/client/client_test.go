package client

import (
	"os"
	"path"
	"testing"

	log "github.com/sirupsen/logrus"
)

func init() {
	// enable debug level log messages
	log.SetLevel(log.DebugLevel)
}

func TestNewWebhook(t *testing.T) {

	c := Webhook{
		QaasHost:     "localhost",
		QaasPort:     443,
		QaasCertFile: path.Join(os.Getenv("GOPATH"), "src/github.com/Donders-Institute/hpc-qaas/test/cert/TestServer.crt"),
	}

	script := path.Join(os.Getenv("GOPATH"), "src/github.com/Donders-Institute/hpc-qaas/test/data/qsub.sh")

	url, err := c.New(script)

	if err != nil {
		t.Errorf("test failed: %+v\n", err)
	} else {
		t.Logf("webhook url: %s\n", url.String())
	}
}

func TestListWebhook(t *testing.T) {
	c := Webhook{
		QaasHost:     "localhost",
		QaasPort:     443,
		QaasCertFile: path.Join(os.Getenv("GOPATH"), "src/github.com/Donders-Institute/hpc-qaas/test/cert/TestServer.crt"),
	}

	chanWebhook, err := c.List()
	if err != nil {
		t.Errorf("test failed: %+v\n", err)
	} else {
		for webhook := range chanWebhook {
			t.Logf("webhook: %+v\n", webhook)
		}
	}
}
