package server

import (
	"bytes"
	"fmt"
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

func TestConfigurationAddHandler(t *testing.T) {
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
				Hash:        "550e8400-e29b-41d4-a716-446655440001",
				Groupname:   "groupname",
				Username:    "username",
				Description: "description",
			},
			testData: `{"hash": "550e8400-e29b-41d4-a716-446655440001", "groupname": "groupname", "username": "username", "description": "description"}`,
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
				Hash:        "550e8400-e29b-41d4-a716-446655440001",
				Groupname:   "groupname",
				Username:    "username",
				Description: "description",
			},
			testData: `{"hash": "550e8400-e29b-41d4-a716-446655440001", "groupname": "groupname", "username": "username", "description": "description"}`,
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
				Hash:        "550e8400-e29b-41d4-a716-446655440001",
				Groupname:   "groupname",
				Username:    "username",
				Description: "description",
			},
			testData: `{"hash": "550e8400-e29b-41d4-a716-446655440001", "groupname": "groupname", "username": "username", "description": "description"}`,
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
			QaasInternalPort:   "443",
			QaasExternalPort:   "5111",
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
			sqlmock.NewRows([]string{"id", "hash", "groupname", "username", "description", "created"}).
				AddRow(1,
					c.configuration.Hash,
					c.configuration.Groupname,
					c.configuration.Username,
					c.configuration.Description,
					"2019-03-11T19:44:44+01:00")

			mock.ExpectBegin()
			mock.ExpectExec("INSERT INTO qaas").
				WithArgs(c.configuration.Hash,
					c.configuration.Groupname,
					c.configuration.Username,
					c.configuration.Description,
					"2019-03-11T19:44:44+01:00").
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
		}

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.ConfigurationAddHandler)

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

func TestConfigurationInfoHandler(t *testing.T) {
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
			method:    "GET",
			configURL: "/configuration/550e8400-e29b-41d4-a716-446655440001",
			configuration: ConfigurationRequest{
				Hash:        "550e8400-e29b-41d4-a716-446655440001",
				Groupname:   "groupname",
				Username:    "username",
				Description: "description",
			},
			testData: `{"hash": "550e8400-e29b-41d4-a716-446655440001", "groupname": "groupname", "username": "username", description": "description"}`,
			headerInfo: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			expectedStatus: 200,
			expectedString: `{"webhook":{"hash":"550e8400-e29b-41d4-a716-446655440001","groupname":"groupname","username":"username",description":"description"}}`,
			expectedResult: true, // No error
		},
		{
			method:    "GET",
			configURL: "/configuration/nonexisting",
			configuration: ConfigurationRequest{
				Hash:        "550e8400-e29b-41d4-a716-446655440001",
				Groupname:   "groupname",
				Username:    "username",
				Description: "description",
			},
			testData: `{"hash": "550e8400-e29b-41d4-a716-446655440001", "groupname": "groupname", "username": "username", description": "description"}`,
			headerInfo: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			expectedStatus: 404,
			expectedString: `Error 404 - Not found: invalid URL path '/configuration/nonexisting'`,
			expectedResult: false, // No error
		},
		{
			method:    "POST",
			configURL: "/configuration/550e8400-e29b-41d4-a716-446655440001",
			configuration: ConfigurationRequest{
				Hash:        "550e8400-e29b-41d4-a716-446655440001",
				Groupname:   "groupname",
				Username:    "username",
				Description: "description",
			},
			testData: `{"hash": "550e8400-e29b-41d4-a716-446655440001", "groupname": "groupname", "username": "username", description": "description"}`,
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
			QaasInternalPort:   "443",
			QaasExternalPort:   "5111",
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
			expectedRows := sqlmock.NewRows([]string{"id", "hash", "groupname", "username", "description", "created"}).
				AddRow(1,
					c.configuration.Hash,
					c.configuration.Groupname,
					c.configuration.Username,
					c.configuration.Description,
					"2019-03-11T19:44:44+01:00",
				)
			mock.ExpectQuery("^SELECT id, hash, groupname, username, description, created FROM qaas").
				WithArgs(c.configuration.Hash).
				WillReturnRows(expectedRows)
		}

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.ConfigurationInfoHandler)

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

func TestConfigurationListHandler(t *testing.T) {
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
			method:    "GET",
			configURL: "/configuration",
			configuration: ConfigurationRequest{
				Hash:        "",
				Groupname:   "groupname",
				Username:    "username",
				Description: "",
			},
			testData: `{"hash": "", "groupname": "groupname", "username": "username", description": ""}`,
			headerInfo: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			expectedStatus: 200,
			expectedString: `{"webhooks":[{"hash":"550e8400-e29b-41d4-a716-446655440001","groupname":"groupname","username":"username","description":"","created":"","url":"https://qaas.dccn.nl:5111/550e8400-e29b-41d4-a716-446655440001"},{"hash":"550e8400-e29b-41d4-a716-446655440002","groupname":"groupname","username":"username","description":"","created":"","url":"https://qaas.dccn.nl:5111/550e8400-e29b-41d4-a716-446655440002"}]}`,
			expectedResult: true, // No error
		},
		{
			method:    "GET",
			configURL: "/configuration/nonexisting",
			configuration: ConfigurationRequest{
				Hash:        "",
				Groupname:   "groupname",
				Username:    "username",
				Description: "",
			},
			testData: `{"hash": "", "groupname": "groupname", "username": "username", "description": ""}`,
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
				Hash:        "",
				Groupname:   "groupname",
				Username:    "username",
				Description: "",
			},
			testData: `{"hash": "", "groupname": "groupname", "username": "username", "description": ""}`,
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
			QaasInternalPort:   "443",
			QaasExternalPort:   "5111",
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
			hash1 := "550e8400-e29b-41d4-a716-446655440001"
			hash2 := "550e8400-e29b-41d4-a716-446655440002"
			expectedRows := sqlmock.NewRows([]string{"id", "hash", "groupname", "username", "description", "created"}).
				AddRow(1,
					hash1,
					c.configuration.Groupname,
					c.configuration.Username,
					c.configuration.Description,
					"2019-03-11T19:44:44+01:00").
				AddRow(2,
					hash2,
					c.configuration.Groupname,
					c.configuration.Username,
					c.configuration.Description,
					"2019-03-11T19:45:44+01:00")
			mock.ExpectQuery("^SELECT id, hash, groupname, username, description, created FROM qaas").
				WithArgs(c.configuration.Groupname, c.configuration.Username).
				WillReturnRows(expectedRows)
		}

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.ConfigurationListHandler)

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

func TestConfigurationDeleteHandler(t *testing.T) {
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
			method:    "DELETE",
			configURL: "/configuration/550e8400-e29b-41d4-a716-446655440001",
			configuration: ConfigurationRequest{
				Hash:        "550e8400-e29b-41d4-a716-446655440001",
				Groupname:   "groupname",
				Username:    "username",
				Description: "description",
			},
			testData: `{"hash": "550e8400-e29b-41d4-a716-446655440001", "groupname": "groupname", "username": "username"}`,
			headerInfo: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			expectedStatus: 200,
			expectedString: `{"webhook":"550e8400-e29b-41d4-a716-446655440001"}`,
			expectedResult: true, // No error
		},
		{
			method:    "DELETE",
			configURL: "/configuration/nonexisting",
			configuration: ConfigurationRequest{
				Hash:        "550e8400-e29b-41d4-a716-446655440001",
				Groupname:   "groupname",
				Username:    "username",
				Description: "description",
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
			configURL: "/configuration/550e8400-e29b-41d4-a716-446655440001",
			configuration: ConfigurationRequest{
				Hash:        "550e8400-e29b-41d4-a716-446655440001",
				Groupname:   "groupname",
				Username:    "username",
				Description: "description",
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
			QaasInternalPort:   "443",
			QaasExternalPort:   "5111",
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
			hash1 := c.configuration.Hash
			hash2 := "550e8400-e29b-41d4-a716-446655440002"
			sqlmock.NewRows([]string{"id", "hash", "groupname", "username", "description", "created"}).
				AddRow(1,
					c.configuration.Hash,
					c.configuration.Groupname,
					c.configuration.Username,
					c.configuration.Description,
					"2019-03-11T19:44:44+01:00").
				AddRow(2,
					hash2,
					c.configuration.Groupname,
					c.configuration.Username,
					c.configuration.Description,
					"2019-03-11T19:45:44+01:00")

			mock.ExpectBegin()
			mock.ExpectExec("DELETE FROM qaas").
				WithArgs(hash1, c.configuration.Groupname, c.configuration.Username).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
		}

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.ConfigurationDeleteHandler)

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
