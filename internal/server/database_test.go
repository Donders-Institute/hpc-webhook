package server

import (
	"reflect"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestAddRow(t *testing.T) {
	configuration := ConfigurationRequest{
		Hash:      "550e8400-e29b-41d4-a716-446655440001",
		Groupname: "dccngroup",
		Username:  "dccnuser",
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO qaas").WithArgs(configuration.Hash, configuration.Groupname, configuration.Username).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err = addRow(db, configuration.Hash, configuration.Groupname, configuration.Username); err != nil {
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
	sqlmock.NewRows([]string{"id", "hash", "groupname", "username"}).
		AddRow(1, hash1, expectedGroupname, expectedUsername).
		AddRow(2, hash2, expectedGroupname, expectedUsername)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM qaas").
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
	expectedRows := sqlmock.NewRows([]string{"id", "hash", "groupname", "username"}).AddRow(1, hash, expectedGroupname, expectedUsername)
	mock.ExpectQuery("^SELECT id, hash, groupname, username FROM qaas").WithArgs(hash).WillReturnRows(expectedRows)

	listExpected := []Item{
		{
			ID:        1,
			Hash:      hash,
			Groupname: expectedGroupname,
			Username:  expectedUsername,
		},
	}

	list, err = getRow(db, hash)
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

	hash2 := "550e8400-e29b-41d4-a716-446655440002"
	expectedGroupname2 := "dccngroup"
	expectedUsername2 := "dccnuser"

	expectedRows := sqlmock.NewRows([]string{"id", "hash", "groupname", "username"}).AddRow(1, hash1, expectedGroupname1, expectedUsername1).AddRow(2, hash2, expectedGroupname2, expectedUsername2)

	mock.ExpectQuery("^SELECT id, hash, groupname, username FROM qaas").WithArgs(expectedGroupname1, expectedUsername1).WillReturnRows(expectedRows)

	listExpected := []Item{
		{
			ID:        1,
			Hash:      hash1,
			Groupname: expectedGroupname1,
			Username:  expectedUsername1,
		},
		{
			ID:        2,
			Hash:      hash2,
			Groupname: expectedGroupname2,
			Username:  expectedUsername2,
		},
	}

	list, err = getListRows(db, expectedGroupname1, expectedUsername2)
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
