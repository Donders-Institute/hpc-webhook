package server

import (
	"fmt"
	"reflect"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestAddRow(t *testing.T) {
	configuration := ConfigurationRequest{
		Hash:        "550e8400-e29b-41d4-a716-446655440001",
		Groupname:   "dccngroup",
		Username:    "dccnuser",
		Description: "description",
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO hpc_webhook").WithArgs(configuration.Hash,
		configuration.Groupname,
		configuration.Username,
		configuration.Description,
		"2019-03-11 10:10:00").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err = addRow(db,
		configuration.Hash,
		configuration.Groupname,
		configuration.Username,
		configuration.Description,
		"2019-03-11 10:10:00"); err != nil {
		t.Errorf("error was not expected while adding row: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteRow(t *testing.T) {
	var err error
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	hash1 := "550e8400-e29b-41d4-a716-446655440001"
	hash2 := "550e8400-e29b-41d4-a716-446655440002"

	expectedGroupname := "dccngroup"
	expectedUsername := "dccnuser"
	sqlmock.NewRows([]string{"id", "hash", "groupname", "username", "description", "created"}).
		AddRow(1, hash1, expectedGroupname, expectedUsername, "This is script 1", "2019-03-11 10:10:00").
		AddRow(2, hash2, expectedGroupname, expectedUsername, "This is script 2", "2019-03-11 10:20:00")

	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM hpc_webhook").
		WithArgs(hash2, expectedGroupname, expectedUsername).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := deleteRow(db, hash2, expectedGroupname, expectedUsername); err != nil {
		t.Errorf("error was not expected while getting row: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetRowHashOnly(t *testing.T) {
	var list []Item
	var err error
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	hash := "550e8400-e29b-41d4-a716-446655440001"

	expectedGroupname := "dccngroup"
	expectedUsername := "dccnuser"
	expectedDescription := "This is script 1"
	expectedCreated := "2019-03-11 10:10:00"
	expectedRows := sqlmock.NewRows([]string{"id", "hash", "groupname", "username", "description", "created"}).
		AddRow(1, hash, expectedGroupname, expectedUsername, expectedDescription, expectedCreated)

	mock.ExpectQuery("^SELECT id, hash, groupname, username, description, created FROM hpc_webhook WHERE").
		WithArgs(hash).
		WillReturnRows(expectedRows)

	hpcWebhookHost := "hpc-webhook.dccn.nl"
	hpcWebhookExternalPort := "443"
	listExpected := []Item{
		{
			ID:          1,
			Hash:        hash,
			Groupname:   expectedGroupname,
			Username:    expectedUsername,
			Description: expectedDescription,
			Created:     expectedCreated,
			URL:         fmt.Sprintf("https://%s:%s%s/%s", hpcWebhookHost, hpcWebhookExternalPort, WebhookPath, hash),
		},
	}

	list, err = getRowHashOnly(db, hpcWebhookHost, hpcWebhookExternalPort, hash)
	if err != nil {
		t.Errorf("error was not expected while getting row: %s", err)
	}

	if !reflect.DeepEqual(list, listExpected) {
		t.Errorf("Lists are not equal: found length %d, but has %d", len(list), len(listExpected))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetRow(t *testing.T) {
	var list []Item
	var err error
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	hash := "550e8400-e29b-41d4-a716-446655440001"

	expectedGroupname := "dccngroup"
	expectedUsername := "dccnuser"
	expectedDescription := "This is script 1"
	expectedCreated := "2019-03-11 10:10:00"
	expectedRows := sqlmock.NewRows([]string{"id", "hash", "groupname", "username", "description", "created"}).
		AddRow(1, hash, expectedGroupname, expectedUsername, expectedDescription, expectedCreated)

	mock.ExpectQuery("^SELECT id, hash, groupname, username, description, created FROM hpc_webhook WHERE").
		WithArgs(hash, expectedGroupname, expectedUsername).
		WillReturnRows(expectedRows)

	hpcWebhookHost := "hpc-webhook.dccn.nl"
	hpcWebhookExternalPort := "443"
	listExpected := []Item{
		{
			ID:          1,
			Hash:        hash,
			Groupname:   expectedGroupname,
			Username:    expectedUsername,
			Description: expectedDescription,
			Created:     expectedCreated,
			URL:         fmt.Sprintf("https://%s:%s%s/%s", hpcWebhookHost, hpcWebhookExternalPort, WebhookPath, hash),
		},
	}

	list, err = getRow(db, hpcWebhookHost, hpcWebhookExternalPort, hash, expectedGroupname, expectedUsername)
	if err != nil {
		t.Errorf("error was not expected while getting row: %s", err)
	}

	if !reflect.DeepEqual(list, listExpected) {
		t.Errorf("Lists are not equal: found length %d, but has %d", len(list), len(listExpected))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetListRows(t *testing.T) {
	var list []Item
	var err error
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	hash1 := "550e8400-e29b-41d4-a716-446655440001"
	expectedGroupname1 := "dccngroup"
	expectedUsername1 := "dccnuser"
	expectedDescription1 := "This is test1"
	expectedCreated1 := "2019-03-11 10:10:00"

	hash2 := "550e8400-e29b-41d4-a716-446655440002"
	expectedGroupname2 := "dccngroup"
	expectedUsername2 := "dccnuser"
	expectedDescription2 := "This is test2"
	expectedCreated2 := "2019-03-11 11:11:00"

	expectedRows := sqlmock.NewRows([]string{"id", "hash", "groupname", "username", "description", "created"}).
		AddRow(1, hash1, expectedGroupname1, expectedUsername1, expectedDescription1, expectedCreated1).
		AddRow(2, hash2, expectedGroupname2, expectedUsername2, expectedDescription2, expectedCreated2)

	mock.ExpectQuery("^SELECT id, hash, groupname, username, description, created FROM hpc_webhook").
		WithArgs(expectedGroupname1, expectedUsername1).
		WillReturnRows(expectedRows)

	hpcWebhookHost := "hpc-webhook.dccn.nl"
	hpcWebhookExternalPort := "443"
	listExpected := []Item{
		{
			ID:          1,
			Hash:        hash1,
			Groupname:   expectedGroupname1,
			Username:    expectedUsername1,
			Description: expectedDescription1,
			Created:     expectedCreated1,
			URL:         fmt.Sprintf("https://%s:%s%s/%s", hpcWebhookHost, hpcWebhookExternalPort, WebhookPath, hash1),
		},
		{
			ID:          2,
			Hash:        hash2,
			Groupname:   expectedGroupname2,
			Username:    expectedUsername2,
			Description: expectedDescription2,
			Created:     expectedCreated2,
			URL:         fmt.Sprintf("https://%s:%s%s/%s", hpcWebhookHost, hpcWebhookExternalPort, WebhookPath, hash2),
		},
	}

	list, err = getListRows(db, hpcWebhookHost, hpcWebhookExternalPort, expectedGroupname1, expectedUsername1)
	if err != nil {
		t.Errorf("error was not expected while getting row: %s", err)
	}

	if !reflect.DeepEqual(list, listExpected) {
		t.Errorf("Lists are not equal: found length %d, but has %d", len(list), len(listExpected))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
