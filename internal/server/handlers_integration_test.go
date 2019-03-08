package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

type testConfiguration struct {
	homeDir            string
	dataDir            string
	keyDir             string
	privateKeyFilename string
	publicKeyFilename  string
}

func setupTestCase(testConfig testConfiguration) error {
	err := os.MkdirAll(testConfig.homeDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error %s when creating %s dir", err, testConfig.homeDir)
	}
	err = os.MkdirAll(testConfig.dataDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error %s when creating %s dir", err, testConfig.dataDir)
	}
	err = os.MkdirAll(testConfig.keyDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error %s when creating %s dir", err, testConfig.keyDir)
	}
	err = generateKeyPair(testConfig.privateKeyFilename, testConfig.publicKeyFilename)
	if err != nil {
		return err
	}
	return err
}

func teardownTestCase(testConfig testConfiguration) error {
	err := os.RemoveAll(testConfig.homeDir)
	if err != nil {
		return fmt.Errorf("error %s when removing %s dir", err, testConfig.homeDir)
	}
	err = os.RemoveAll(testConfig.dataDir)
	if err != nil {
		return fmt.Errorf("error %s when removing %s dir", err, testConfig.dataDir)
	}
	err = os.RemoveAll(testConfig.keyDir)
	if err != nil {
		return fmt.Errorf("error %s when removing %s dir", err, testConfig.keyDir)
	}
	return err
}

func TestConfigurationHandler(t *testing.T) {
	cases := []struct {
		method         string
		configURL      string
		configuration  ConfigurationRequest
		testData       string
		headerInfo     map[string]string
		expectedStatus int
		expectedString string
		expectedResult bool
	}{
		{
			method:    "PUT",
			configURL: "/configuration",
			configuration: ConfigurationRequest{
				Hash:      "550e8400-e29b-41d4-a716-446655440001",
				Groupname: "groupname",
				Username:  "username",
			},
			testData: `{"hash": "550e8400-e29b-41d4-a716-446655440001", "groupname": "groupname", "username": "username"}`,
			headerInfo: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			expectedStatus: 200,
			expectedString: `{"webhook":"https://qaas.dccn.nl:443/webhook/550e8400-e29b-41d4-a716-446655440001"}`,
			expectedResult: true, // No error
		},
		{
			method:    "PUT",
			configURL: "/configuration/nonexisting",
			configuration: ConfigurationRequest{
				Hash:      "550e8400-e29b-41d4-a716-446655440001",
				Groupname: "groupname",
				Username:  "username",
			},
			testData: `{"hash": "550e8400-e29b-41d4-a716-446655440001", "groupname": "groupname", "username": "username"}`,
			headerInfo: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			expectedStatus: 404,
			expectedString: `Error 404 - Not found: invalid URL path '/configuration/nonexisting'`,
			expectedResult: false, // No error
		},
		{
			method:    "POST",
			configURL: "/configuration",
			configuration: ConfigurationRequest{
				Hash:      "550e8400-e29b-41d4-a716-446655440001",
				Groupname: "groupname",
				Username:  "username",
			},
			testData: `{"hash": "550e8400-e29b-41d4-a716-446655440001", "groupname": "groupname", "username": "username"}`,
			headerInfo: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			expectedStatus: 404,
			expectedString: `Error 404 - Not found: invalid method 'POST'`,
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

		api := API{
			DB: db,
			Connector: FakeConnector{
				Description: "fake SSH connection to relay node",
			},
			DataDir:            testConfig.dataDir,
			HomeDir:            testConfig.homeDir,
			RelayNode:          "relaynode.dccn.nl",
			QaasHost:           "qaas.dccn.nl",
			QaasPort:           "5111",
			PrivateKeyFilename: testConfig.publicKeyFilename,
			PublicKeyFilename:  testConfig.privateKeyFilename,
		}

		app := &api

		// Obtain the test data
		b := bytes.NewBuffer([]byte(c.testData))

		// Make a new HTTP POST request with this body
		req, err := http.NewRequest(c.method, c.configURL, b)
		if err != nil {
			t.Fatal(err)
		}

		// Modify the header
		for key, value := range c.headerInfo {
			req.Header.Set(key, value)
		}

		if c.expectedResult {
			mock.ExpectBegin()
			mock.ExpectExec("INSERT INTO qaas").WithArgs(c.configuration.Hash, c.configuration.Groupname, c.configuration.Username).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
		}

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.ConfigurationHandler)

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
		file, err := os.Open(c.testDataFilename)
		if err != nil {
			t.Errorf("Expected no error, but got '%+v'", err.Error())
		}
		defer file.Close()

		body, err := ioutil.ReadAll(file)
		if err != nil {
			t.Errorf("Test data unavailable: %v", err)
		}
		b := bytes.NewBuffer(body)

		// Make a new HTTP POST request with this body
		req, err := http.NewRequest(c.method, c.payloadURL, b)
		if err != nil {
			t.Fatal(err)
		}

		// Modify the header
		for key, value := range c.headerInfo {
			req.Header.Set(key, value)
		}

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
