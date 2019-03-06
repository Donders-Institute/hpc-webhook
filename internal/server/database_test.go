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

	// now we execute our method
	if err = addRow(db, configuration.Hash, configuration.Groupname, configuration.Username); err != nil {
		t.Errorf("error was not expected while adding row: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetRow(t *testing.T) {
	var list []item
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

	listExpected := []item{
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

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
