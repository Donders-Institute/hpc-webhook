package server

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestConfigurationHandlerWebhook(t *testing.T) {
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
				Hash:     "550e8400-e29b-41d4-a716-446655440001",
				Username: "username",
			},
			testData: `{"hash": "550e8400-e29b-41d4-a716-446655440001", "username": "username"}`,
			headerInfo: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			expectedStatus: 200,
			expectedString: `Webhook added successfully`,
			expectedResult: true, // No error
		},
		{
			method:    "PUT",
			configURL: "/configuration/nonexisting",
			configuration: ConfigurationRequest{
				Hash:     "550e8400-e29b-41d4-a716-446655440001",
				Username: "username",
			},
			testData: `{"hash": "550e8400-e29b-41d4-a716-446655440001", "username": "username"}`,
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
				Hash:     "550e8400-e29b-41d4-a716-446655440001",
				Username: "username",
			},
			testData: `{"hash": "550e8400-e29b-41d4-a716-446655440001", "username": "username"}`,
			headerInfo: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			expectedStatus: 404,
			expectedString: `Error 404 - Not found: invalid method 'POST'`,
			expectedResult: false, // Invalid method
		},
	}

	for _, c := range cases {

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		api := API{
			DB:        db,
			RelayNode: "relaynode.dccn.nl",
			QaasHost:  "qaas.dccn.nl",
			QaasPort:  "5111",
		}
		api.SetDataDir("..", "..", "test", "results")
		err = os.MkdirAll(api.DataDir, os.ModePerm)
		if err != nil {
			t.Fatalf("error %s when creating %s dir", err, api.DataDir)
		}
		defer func() {
			err = os.RemoveAll(api.DataDir) // cleanup when done
			if err != nil {
				t.Fatalf("error %s when removing %s dir", err, api.DataDir)
			}
		}()
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
			mock.ExpectExec("INSERT INTO qaas").WithArgs(c.configuration.Hash, c.configuration.Username).WillReturnResult(sqlmock.NewResult(1, 1))
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
			testDataFilename: "../../test/data/example-github-webhook.json",
			headerInfo: map[string]string{
				"Content-Type":      "application/json; charset=utf-8",
				"x-hub-signature":   "someValue",
				"x-github-event":    "someValue",
				"x-github-delivery": "someValue",
			},
			expectedStatus: 200,
			expectedString: `{"webhook":"https://qaas.dccn.nl:5111/webhook/550e8400-e29b-41d4-a716-446655440001"}`,
			expectedResult: true, // No error
		},
		{
			method:           "POST",
			payloadURL:       "/webhook/550e8400-e29b-41d4-a716-446655440002",
			hash:             "550e8400-e29b-41d4-a716-446655440002",
			testDataFilename: "../../test/data/example-ifttt-webhook.json",
			headerInfo: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			expectedStatus: 200,
			expectedString: `{"webhook":"https://qaas.dccn.nl:5111/webhook/550e8400-e29b-41d4-a716-446655440002"}`,
			expectedResult: true, // No error
		},
		{
			method:           "POST",
			payloadURL:       "/webhook/550e8400-e29b-41d4-a716-446655440003",
			hash:             "550e8400-e29b-41d4-a716-446655440003",
			testDataFilename: "../../test/data/example-zapier-webhook.json",
			headerInfo: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			expectedStatus: 200,
			expectedString: `{"webhook":"https://qaas.dccn.nl:5111/webhook/550e8400-e29b-41d4-a716-446655440003"}`,
			expectedResult: true, // No error
		},
		{
			method:           "POST",
			payloadURL:       "/webhook/550e8400-e29b-41d4-a716",
			hash:             "550e8400-e29b-41d4-a716",
			testDataFilename: "../../test/data/example-zapier-webhook.json",
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
			testDataFilename: "../../test/data/example-zapier-webhook.json",
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
			testDataFilename: "../../test/data/example-zapier-webhook.json",
			headerInfo: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			expectedStatus: 404,
			expectedString: `Error 404 - Not found: invalid method 'GET'`,
			expectedResult: false, // Invalid method
		},
	}

	for _, c := range cases {

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		api := API{
			DB:        db,
			RelayNode: "relaynode.dccn.nl",
			QaasHost:  "qaas.dccn.nl",
			QaasPort:  "5111",
		}
		api.SetDataDir("..", "..", "test", "results")
		err = os.MkdirAll(api.DataDir, os.ModePerm)
		if err != nil {
			t.Fatalf("error %s when creating %s dir", err, api.DataDir)
		}
		defer func() {
			err = os.RemoveAll(api.DataDir) // cleanup when done
			if err != nil {
				t.Fatalf("error %s when removing %s dir", err, api.DataDir)
			}
		}()
		app := &api

		// Obtain the body
		file, err := os.Open(c.testDataFilename)
		if err != nil {
			log.Fatal(err)
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
			expectedUsername := "dccnuser"
			expectedRows := sqlmock.NewRows([]string{"id", "hash", "username"}).AddRow(1, c.hash, expectedUsername)
			mock.ExpectQuery("^SELECT id, hash, username FROM qaas").WithArgs(c.hash).WillReturnRows(expectedRows)
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
