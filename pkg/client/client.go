package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"
)

type requestDataQaas struct {
	Username string `json:"username"`
	Hash     string `json:"hash"`
}

type responseDataQaas struct {
}

// WebhookInfo is a data structure containing the information (and/or attributes) of a webhook.
type WebhookInfo struct {
	ID string
}

// Webhook provides client interfaces for managing webhook registry on the QaaS server, using the RESTful interface.
type Webhook struct {
	QaasHost     string
	QaasPort     int
	QaasCertFile string
}

// New provisions a new webhook for QaaS and registry the new webhook at the QaaS server.
func (s *Webhook) New(script string) (*url.URL, error) {

	// check existance of the script, and it's type.
	scriptAbs, err := filepath.Abs(script)
	if err != nil {
		return nil, err
	}
	fi, err := os.Lstat(scriptAbs)
	if err != nil {
		return nil, err
	}
	if !fi.Mode().IsRegular() {
		return nil, fmt.Errorf("not a regular file: %s", script)
	}

	// get current user
	user, err := user.Current()
	if err != nil {
		return nil, err
	}
	id := uuid.New().String()
	workdir := path.Join(user.HomeDir, ".qaas", id)

	if err := os.MkdirAll(workdir, 0700); err != nil {
		return nil, err
	}

	// provision necessary directory
	// - write path to the script.sh file
	f, err := os.Create(path.Join(workdir, "script.sh"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if _, err := f.WriteString(fmt.Sprintf("%s\n", scriptAbs)); err != nil {
		return nil, err
	}

	// call QaaS to register the webhook
	myURL := fmt.Sprintf("https://%s:%d/webhooks/%s", s.QaasHost, s.QaasPort, id)
	var respData responseDataQaas
	httpCode, err := s.putJSON(myURL, requestDataQaas{Username: user.Name, Hash: id}, respData)

	log.Debugf("response data: %+v", respData)

	if err != nil || httpCode != 200 {
		return nil, fmt.Errorf("error registering webhook on QaaS server: +%v (HTTP CODE: %d)", err, httpCode)
	}

	return url.Parse(myURL)
}

// List retrieves a list of webhooks of the current user.
func (s *Webhook) List() ([]WebhookInfo, error) {

	var webhooks []WebhookInfo

	// get current user
	user, err := user.Current()
	if err != nil {
		return webhooks, err
	}

	// add names of the items under $HOME/.gass into the list if:
	//
	// - the item is a directory
	// - the name of the item can be pased by uuid.Parse() function
	if items, err := ioutil.ReadDir(path.Join(user.HomeDir, ".qaas")); err == nil {
		for _, f := range items {
			if f.IsDir() {
				if _, err := uuid.Parse(f.Name()); err == nil {
					webhooks = append(webhooks, WebhookInfo{ID: f.Name()})
				}
			}
		}
	}

	return webhooks, nil
}

// putJSON makes a HTTP PUT request with provided JSON data.
func (s *Webhook) putJSON(url string, request interface{}, response interface{}) (int, error) {

	data, err := json.Marshal(request)
	if err != nil {
		return 0, err
	}

	log.Debugf("request data: %s", string(data))

	c := s.newHTTPSClient()
	req, err := http.NewRequest("PUT", url, bytes.NewReader(data))
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	if err != nil {
		return 0, err
	}

	// make HTTP PUT call
	rsp, err := c.Do(req)
	if err != nil {
		return 0, err
	}
	defer rsp.Body.Close()

	return rsp.StatusCode, json.NewDecoder(rsp.Body).Decode(response)
}

// newHTTPSClient sets up the client instance ready for making HTTPs requests.
func (s *Webhook) newHTTPSClient() *http.Client {

	rootCertPool := x509.NewCertPool()

	if s.QaasCertFile != "" {
		pem, _ := ioutil.ReadFile(s.QaasCertFile)
		rootCertPool.AppendCertsFromPEM(pem)
	}

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 5 * time.Second,
		TLSClientConfig: &tls.Config{
			RootCAs: rootCertPool,
		},
	}

	return &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}
}
