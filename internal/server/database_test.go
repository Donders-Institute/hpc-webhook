package server

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestAddRow(t *testing.T) {
	configuration := Configuration{
		Hash:     "e66d248b67c0442fe2cbad7e248651fd4569ee8ecc72ee5a19b0e55ac1ef4492",
		Username: "dccnuser",
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO qaas").WithArgs(configuration.Hash, configuration.Username).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// now we execute our method
	if err = addRow(db, configuration.Hash, configuration.Username); err != nil {
		t.Errorf("error was not expected while adding row: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
