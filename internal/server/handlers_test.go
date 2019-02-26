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
		configuration  Configuration
		testData       string
		headerInfo     map[string]string
		expectedStatus int
		expectedString string
		expectedResult bool
	}{
		{
			method:    "PUT",
			configURL: "/configuration",
			configuration: Configuration{
				Hash:     "e66d248b67c0442fe2cbad7e248651fd4569ee8ecc72ee5a19b0e55ac1ef4493",
				Username: "username",
			},
			testData: `{"hash": "e66d248b67c0442fe2cbad7e248651fd4569ee8ecc72ee5a19b0e55ac1ef4493", "username": "username"}`,
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
			configuration: Configuration{
				Hash:     "e66d248b67c0442fe2cbad7e248651fd4569ee8ecc72ee5a19b0e55ac1ef4493",
				Username: "username",
			},
			testData: `{"hash": "e66d248b67c0442fe2cbad7e248651fd4569ee8ecc72ee5a19b0e55ac1ef4493", "username": "username"}`,
			headerInfo: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			expectedStatus: 404,
			expectedString: `Error 404 - Not found: invalid URL path`,
			expectedResult: false, // No error
		},
		{
			method:    "POST",
			configURL: "/configuration",
			configuration: Configuration{
				Hash:     "e66d248b67c0442fe2cbad7e248651fd4569ee8ecc72ee5a19b0e55ac1ef4493",
				Username: "username",
			},
			testData: `{"hash": "e66d248b67c0442fe2cbad7e248651fd4569ee8ecc72ee5a19b0e55ac1ef4493", "username": "username"}`,
			headerInfo: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			expectedStatus: 404,
			expectedString: `Error 404 - Not found: invalid method`,
			expectedResult: false, // Invalid method
		},
	}

	for _, c := range cases {

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		app := &API{db}

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
		testDataFilename string
		headerInfo       map[string]string
		expectedStatus   int
		expectedString   string
		expectedResult   bool
	}{
		{
			method:           "POST",
			payloadURL:       "/webhook/e66d248b67c0442fe2cbad7e248651fd4569ee8ecc72ee5a19b0e55ac1ef4492",
			testDataFilename: "../../test/data/example-github-webhook.json",
			headerInfo: map[string]string{
				"Content-Type":      "application/json; charset=utf-8",
				"x-hub-signature":   "someValue",
				"x-github-event":    "someValue",
				"x-github-delivery": "someValue",
			},
			expectedStatus: 200,
			expectedString: `Webhook payload delivered successfully`,
			expectedResult: true, // No error
		},
		{
			method:           "POST",
			payloadURL:       "/webhook/e66d248b67c0442fe2cbad7e248651fd4569ee8ecc72ee5a19b0e55ac1ef4493",
			testDataFilename: "../../test/data/example-ifttt-webhook.json",
			headerInfo: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			expectedStatus: 200,
			expectedString: `Webhook payload delivered successfully`,
			expectedResult: true, // No error
		},
		{
			method:           "POST",
			payloadURL:       "/webhook/e66d248b67c0442fe2cbad7e248651fd4569ee8ecc72ee5a19b0e55ac1ef4494",
			testDataFilename: "../../test/data/example-zapier-webhook.json",
			headerInfo: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			expectedStatus: 200,
			expectedString: `Webhook payload delivered successfully`,
			expectedResult: true, // No error
		},
		{
			method:           "POST",
			payloadURL:       "/webhook/e66d248b67c0442fe2cbad7e248651",
			testDataFilename: "../../test/data/example-zapier-webhook.json",
			headerInfo: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			expectedStatus: 404,
			expectedString: `Error 404 - Not found: invalid URL path`,
			expectedResult: false, // Invalid method
		},
		{
			method:           "POST",
			payloadURL:       "/wwwhook/e66d248b67c0442fe2cbad7e248651fd4569ee8ecc72ee5a19b0e55ac1ef4494",
			testDataFilename: "../../test/data/example-zapier-webhook.json",
			headerInfo: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			expectedStatus: 404,
			expectedString: `Error 404 - Not found: invalid URL path`,
			expectedResult: false, // Invalid method
		},
		{
			method:           "GET",
			payloadURL:       "/webhook/e66d248b67c0442fe2cbad7e248651fd4569ee8ecc72ee5a19b0e55ac1ef4494",
			testDataFilename: "../../test/data/example-zapier-webhook.json",
			headerInfo: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			expectedStatus: 404,
			expectedString: `Error 404 - Not found: invalid method`,
			expectedResult: false, // Invalid method
		},
	}

	for _, c := range cases {

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		app := &API{DB: db}

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

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}
