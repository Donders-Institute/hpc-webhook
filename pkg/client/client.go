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
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/Donders-Institute/hpc-qaas/internal/server"

	log "github.com/sirupsen/logrus"
)

// WebhookConfigInfo is a data structure containing the information (and/or attributes) of a webhook.
type WebhookConfigInfo struct {
	ID           string
	Description  string
	CreationTime string
	Script       string
	WebhookURL   string
}

// TriggerWebhook makes a POST call to the WebhookURL with the given payload.
//
// The payload is specified by `data` as a byte array, and the string `contentType` as the
// HTTP request content type.
//
// For the WebhookURL supporting HTTPS protocol, the provided X509 certificate file `cacert`
// is used for making the connection.
//
// The response body of the POST is returned as a byte array.
func (info *WebhookConfigInfo) TriggerWebhook(payload []byte, contentType string, cacert string) ([]byte, error) {

	if info.ID == "" || info.WebhookURL == "" || info.Script == "" {
		return nil, fmt.Errorf("invalid Webhook: %+v", info)
	}

	myURL, err := url.Parse(info.WebhookURL)
	if err != nil {
		return nil, err
	}

	c := httpsClient(cacert)
	req, err := http.NewRequest("POST", myURL.String(), bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", contentType)

	// make HTTP POST call
	rsp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {
		return nil, fmt.Errorf("fail trigger webhook: %s (%d: %s)", myURL.String(), rsp.StatusCode, rsp.Status)
	}

	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	return rspBody, nil
}

// WebhookConfig provides client interfaces for managing webhook registry on the QaaS server, using the RESTful interface.
type WebhookConfig struct {
	QaasHost     string
	QaasPort     int
	QaasCertFile string
}

// New provisions a new WebhookConfig for QaaS and registry the new webhook at the QaaS server.
func (s *WebhookConfig) New(script string, desc string) (*url.URL, error) {

	// check existence of the script and its type.
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
	cuser, err := user.Current()
	if err != nil {
		return nil, err
	}
	id := uuid.New().String()
	workdir := path.Join(cuser.HomeDir, ".qaas", id)

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
	myURL := url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s:%d", s.QaasHost, s.QaasPort),
		Path:   server.ConfigurationPath,
	}
	var response server.ConfigurationResponse

	cgroup, err := user.LookupGroupId(cuser.Gid)
	if err != nil {
		return nil, err
	}
	httpCode, err := httpPutJSON(
		&myURL,
		s.QaasCertFile,
		&server.ConfigurationRequest{
			Hash:        id,
			Groupname:   cgroup.Name,
			Username:    cuser.Username,
			Description: desc,
		},
		&response)

	log.Debugf("response data: %+v", response)

	if err != nil {
		return nil, fmt.Errorf("error registering webhook on QaaS server: +%v (HTTP CODE: %d)", err, httpCode)
	}

	webhookURL, err := url.Parse(response.Webhook)
	if err != nil {
		return nil, err
	}

	return webhookURL, nil
}

// List retrieves a list of webhooks of the current user.
// The information of webhooks is returned with a channel.
func (s *WebhookConfig) List() (chan WebhookConfigInfo, error) {

	// channel for webhook ids found in local QaaS directory.
	chanWebhookID := make(chan string)

	// channel for webhook information found in the remote QaaS server.
	chanWebhookConfigInfo := make(chan WebhookConfigInfo)

	wg := new(sync.WaitGroup)
	nworkers := 4
	wg.Add(nworkers)

	for i := 0; i < nworkers; i++ {
		go func() {
			for id := range chanWebhookID {
				info, err := s.GetInfo(id)
				if err != nil {
					// only show error message in debug mode
					log.Debugln(err)
				} else {
					chanWebhookConfigInfo <- info
				}
			}
			wg.Done()
		}()
	}

	// go routine feeding webhook ids to chanWebhookID, and wait for all local webhook ids are visited to get webhookInfo
	go func() {
		// perform ID fetching from user's home directory only when the current user can be determined.
		if cuser, err := user.Current(); err == nil {
			// add names of the items under $HOME/.gass into the list if:
			//
			// - the item is a directory
			// - the name of the item can be passed by uuid.Parse() function
			if items, err := ioutil.ReadDir(path.Join(cuser.HomeDir, ".qaas")); err == nil {
				for _, f := range items {
					if !f.IsDir() {
						continue
					}
					if _, err := uuid.Parse(f.Name()); err == nil {
						chanWebhookID <- f.Name()
					}
				}
			}
		}
		close(chanWebhookID)

		wg.Wait()
		close(chanWebhookConfigInfo)
	}()

	return chanWebhookConfigInfo, nil
}

// GetInfo retrieves information of a single Webhook configuration referred by the hash id.
func (s *WebhookConfig) GetInfo(id string) (WebhookConfigInfo, error) {

	info := WebhookConfigInfo{ID: id}

	// get current user
	cuser, err := user.Current()
	if err != nil {
		return info, err
	}

	// get current user's primary group
	cgroup, err := user.LookupGroupId(cuser.Gid)
	if err != nil {
		return info, err
	}

	myURL := url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s:%d", s.QaasHost, s.QaasPort),
		Path:   path.Join(server.ConfigurationPath, id),
	}
	var response server.ConfigurationInfoResponse

	httpCode, err := httpGetJSON(
		&myURL,
		s.QaasCertFile,
		&server.ConfigurationRequest{
			Hash:        id,
			Groupname:   cgroup.Name,
			Username:    cuser.Username,
			Description: "",
		},
		&response)

	log.Debugf("response data: %+v", response)

	if err != nil {
		return info, fmt.Errorf("error retrieving webhook info from the QaaS server: %+v (HTTP CODE: %d)", err, httpCode)
	}

	if id != response.Webhook.Hash {
		return info, fmt.Errorf("expect webhook id: %s, server returns id: %s", id, response.Webhook.Hash)
	}

	info.Description = response.Webhook.Description
	info.CreationTime = response.Webhook.Created
	info.WebhookURL = response.Webhook.URL

	// read local script from the webhook's working directory
	if script, err := ioutil.ReadFile(path.Join(cuser.HomeDir, ".qaas", id, "script.sh")); err != nil {
		log.Errorf("cannot locate script of webhook: %s\n", id)
	} else {
		// remove tailing "\n"
		info.Script = strings.TrimSuffix(string(script), "\n")
	}

	return info, nil
}

// Delete removes a webhook with the given id.
//
// The deletion maily removes webhook registry from QaaS server.
// If removeDir is true, the local webhook working directory is removed when the webhook is unregistered from the QaaS server.
func (s *WebhookConfig) Delete(id string, removeDir bool) error {

	// check if there is a webhook directory in user's .qaas directory.
	cuser, err := user.Current()
	if err != nil {
		return err
	}
	workdir := path.Join(cuser.HomeDir, ".qaas", id)

	w, err := os.Lstat(workdir)
	if err != nil {
		return err
	}
	if !w.IsDir() {
		return fmt.Errorf("not a directory: %s", workdir)
	}

	// check if we can get the given webhook from the server.
	if _, err := s.GetInfo(id); err != nil {
		return fmt.Errorf("webhook not found: %s", err)
	}

	// make DELETE call to the server and receive response.
	myURL := url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s:%d", s.QaasHost, s.QaasPort),
		Path:   path.Join(server.ConfigurationPath, id),
	}

	cgroup, err := user.LookupGroupId(cuser.Gid)
	if err != nil {
		return err
	}
	httpCode, err := httpDelete(
		&myURL,
		s.QaasCertFile,
		&server.ConfigurationRequest{
			Hash:        id,
			Groupname:   cgroup.Name,
			Username:    cuser.Username,
			Description: "",
		},
	)

	if err != nil {
		return fmt.Errorf("fail to delete webhook %s: %+v (HTTP CODE: %d)", id, err, httpCode)
	}

	// remove webhook folder conditionally
	if removeDir {
		return os.RemoveAll(workdir)
	}

	return nil
}

// httpPutJSON makes a HTTP PUT request with provided JSON data.
func httpPutJSON(url *url.URL, cacert string, request interface{}, response interface{}) (int, error) {

	data, err := json.Marshal(request)
	if err != nil {
		return 0, err
	}

	log.Debugf("request data: %s", string(data))

	c := httpsClient(cacert)
	req, err := http.NewRequest("PUT", url.String(), bytes.NewReader(data))
	if err != nil {
		return 0, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")

	// make HTTP PUT call
	rsp, err := c.Do(req)
	if err != nil {
		return 0, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {
		return rsp.StatusCode, fmt.Errorf("%s", rsp.Status)
	}

	return rsp.StatusCode, json.NewDecoder(rsp.Body).Decode(response)
}

// httpGetJSON makes a HTTP GET request to the given url and returns unmarshals JSON response.
func httpGetJSON(url *url.URL, cacert string, request interface{}, response interface{}) (int, error) {

	c := httpsClient(cacert)

	var req *http.Request

	if request != nil {
		// with JSON request body in the GET call
		data, err := json.Marshal(request)
		if err != nil {
			return 0, err
		}
		log.Debugf("request data: %s", string(data))
		req, err = http.NewRequest("GET", url.String(), bytes.NewReader(data))
		if err != nil {
			return 0, err
		}
		req.Header.Set("content-type", "application/json")
	} else {
		// without request body in the GET call
		var err error
		req, err = http.NewRequest("GET", url.String(), nil)
		if err != nil {
			return 0, err
		}
	}

	req.Header.Set("content-type", "application/json")

	// make HTTP GET call
	rsp, err := c.Do(req)
	if err != nil {
		return 0, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {
		return rsp.StatusCode, fmt.Errorf("%s", rsp.Status)
	}

	return rsp.StatusCode, json.NewDecoder(rsp.Body).Decode(response)
}

// httpDelete makes a HTTP DELETE request to the given url.
func httpDelete(url *url.URL, cacert string, request interface{}) (int, error) {

	data, err := json.Marshal(request)
	if err != nil {
		return 0, err
	}

	log.Debugf("request data: %s", string(data))

	c := httpsClient(cacert)
	req, err := http.NewRequest("DELETE", url.String(), bytes.NewReader(data))
	if err != nil {
		return 0, err
	}
	req.Header.Set("content-type", "application/json")

	// make HTTP DELETE call
	rsp, err := c.Do(req)
	if err != nil {
		return 0, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {
		return rsp.StatusCode, fmt.Errorf("%s", rsp.Status)
	}

	return rsp.StatusCode, nil
}

// newHTTPSClient sets up the client instance ready for making HTTPs requests.
func httpsClient(cacert string) *http.Client {

	rootCertPool := x509.NewCertPool()

	if cacert != "" {
		pem, _ := ioutil.ReadFile(cacert)
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
