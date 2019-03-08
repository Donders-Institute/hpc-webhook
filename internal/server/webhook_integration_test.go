package server

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func obtainWebhookPayloadBody(testDataFilename string) (*bytes.Buffer, error) {
	file, err := os.Open(testDataFilename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(body), err
}

func TestHandlerWebhook(t *testing.T) {
	cases := []struct {
		method           string
		payloadURL       string
		hash             string
		username         string
		groupname        string
		testDataFilename string
		headerInfo       map[string]string
		expectedStatus   int
		expectedString   string
		expectedResult   bool
	}{
		{
			method:           "POST",
			payloadURL:       "/webhook/550e8400-e29b-41d4-a716-446655440001",
			hash:             "550e8400-e29b-41d4-a716-446655440001",
			username:         "dccnuser",
			groupname:        "dccngroup",
			testDataFilename: path.Join("..", "..", "test", "data", "example-github-webhook.json"),
			headerInfo: map[string]string{
				"Content-Type":      "application/json; charset=utf-8",
				"x-hub-signature":   "someValue",
				"x-github-event":    "someValue",
				"x-github-delivery": "someValue",
			},
			expectedStatus: 200,
			expectedString: "Webhook handled successfully",
			expectedResult: true, // No error
		},
		{
			method:           "POST",
			payloadURL:       "/webhook/550e8400-e29b-41d4-a716-446655440002",
			hash:             "550e8400-e29b-41d4-a716-446655440002",
			username:         "dccnuser",
			groupname:        "dccngroup",
			testDataFilename: path.Join("..", "..", "test", "data", "example-ifttt-webhook.json"),
			headerInfo: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			expectedStatus: 200,
			expectedString: "Webhook handled successfully",
			expectedResult: true, // No error
		},
		{
			method:           "POST",
			payloadURL:       "/webhook/550e8400-e29b-41d4-a716-446655440003",
			hash:             "550e8400-e29b-41d4-a716-446655440003",
			username:         "dccnuser",
			groupname:        "dccngroup",
			testDataFilename: path.Join("..", "..", "test", "data", "example-zapier-webhook.json"),
			headerInfo: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			expectedStatus: 200,
			expectedString: "Webhook handled successfully",
			expectedResult: true, // No error
		},
		{
			method:           "POST",
			payloadURL:       "/webhook/550e8400-e29b-41d4-a716",
			hash:             "550e8400-e29b-41d4-a716",
			username:         "dccnuser",
			groupname:        "dccngroup",
			testDataFilename: path.Join("..", "..", "test", "data", "example-zapier-webhook.json"),
			headerInfo: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			expectedStatus: 404,
			expectedString: `Error 404 - Not found: invalid URL path '/webhook/550e8400-e29b-41d4-a716'`,
			expectedResult: false, // Invalid webhook id
		},
		{
			method:           "POST",
			payloadURL:       "/wwwhook/550e8400-e29b-41d4-a716-446655440001",
			hash:             "550e8400-e29b-41d4-a716-446655440001",
			username:         "dccnuser",
			groupname:        "dccngroup",
			testDataFilename: path.Join("..", "..", "test", "data", "example-zapier-webhook.json"),
			headerInfo: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			expectedStatus: 404,
			expectedString: `Error 404 - Not found: invalid URL path '/wwwhook/550e8400-e29b-41d4-a716-446655440001'`,
			expectedResult: false, // Invalid URL path
		},
		{
			method:           "GET",
			payloadURL:       "/webhook/550e8400-e29b-41d4-a716-446655440001",
			hash:             "550e8400-e29b-41d4-a716-446655440001",
			username:         "dccnuser",
			groupname:        "dccngroup",
			testDataFilename: path.Join("..", "..", "test", "data", "example-zapier-webhook.json"),
			headerInfo: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			expectedStatus: 404,
			expectedString: `Error 404 - Not found: invalid method 'GET'`,
			expectedResult: false, // Invalid method
		},
	}

	keyDir := path.Join("..", "..", "test", "results", "keys")
	testConfig := testConfiguration{
		homeDir:            path.Join("..", "..", "test", "results", "home"),
		dataDir:            path.Join("..", "..", "test", "results", "data"),
		keyDir:             keyDir,
		privateKeyFilename: path.Join(keyDir, "qaas"),
		publicKeyFilename:  path.Join(keyDir, "qaas.pub"),
	}

	err := setupTestCase(testConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := teardownTestCase(testConfig); err != nil {
			t.Fatal(err)
		}
	}()

	for _, c := range cases {

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		// Configure the application
		api := API{
			DB: db,
			Connector: FakeConnector{
				Description: "fake SSH connection to relay node",
			},
			DataDir:                   testConfig.dataDir,
			HomeDir:                   testConfig.homeDir,
			RelayNode:                 "relaynode.dccn.nl",
			RelayNodeTestUser:         c.username,
			RelayNodeTestUserPassword: "somepassword",
			QaasHost:                  "qaas.dccn.nl",
			QaasPort:                  "5111",
			PrivateKeyFilename:        testConfig.privateKeyFilename,
			PublicKeyFilename:         testConfig.publicKeyFilename,
		}

		app := &api

		// Create the user script file
		userScriptDir := path.Join(api.HomeDir, c.groupname, c.username, ".qaas", c.hash)
		userScriptPathFilename := path.Join(userScriptDir, "script.sh")
		err = os.MkdirAll(userScriptDir, os.ModePerm)
		if err != nil {
			t.Errorf("Error writing user script dir")
		}
		err = ioutil.WriteFile(userScriptPathFilename, []byte("test.sh"), 0644)
		if err != nil {
			t.Errorf("Error writing script.sh")
		}

		// Obtain the body
		b, err := obtainWebhookPayloadBody(c.testDataFilename)
		if err != nil {
			t.Fatal(err)
		}

		// Make a new HTTP POST request with this body
		req, err := http.NewRequest(c.method, c.payloadURL, b)
		if err != nil {
			t.Fatal(err)
		}

		// Modify the header
		for key, value := range c.headerInfo {
			req.Header.Set(key, value)
		}

		// Set the query that is expected to be executed
		if c.expectedResult {
			expectedUsername := c.username
			expectedGroupname := c.groupname
			expectedRows := sqlmock.NewRows([]string{"id", "hash", "groupname", "username"}).AddRow(1, c.hash, expectedGroupname, expectedUsername)
			mock.ExpectQuery("^SELECT id, hash, groupname, username FROM qaas").WithArgs(c.hash).WillReturnRows(expectedRows)
		}

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.WebhookHandler)

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != c.expectedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v", status, c.expectedStatus)
			return
		}

		// Check the expected string
		if rr.Body.String() != c.expectedString {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), c.expectedString)
			return
		}

		if c.expectedResult {
			// we make sure that all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		}
	}
}
