package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"go_project_template/internal/user/repository"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func NewDBMock() (*sql.DB, sqlmock.Sqlmock) {
	// creates sqlmock database connection and a mock to manage expectations
	db, mock, err := sqlmock.New()

	if err != nil {
		panic(err)
	}

	return db, mock
}

func TestGetUser(t *testing.T) {
	// creates sqlmock database connection and a mock to manage expectations
	db, mock := NewDBMock()
	defer db.Close()

	// Setup Repository
	Repository := repository.NewUserRepository(db)

	// Expected query statement
	queryStatement := `SELECT * FROM public.users`

	// Filling mock database
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(int64(1), "Rizky")
	rows.AddRow(int64(2), "Ardi")
	rows.AddRow(int64(3), "Maulana")

	// Set expected outcome
	mock.ExpectQuery(regexp.QuoteMeta(queryStatement)).WillReturnRows(rows)

	// Excecute mock query
	users, err := Repository.GetUsers(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, users)
}

func TestGetUserById(t *testing.T) {
	// creates sqlmock database connection and a mock to manage expectations
	db, mock := NewDBMock()
	defer db.Close()
	// Setup Repository
	Repository := repository.NewUserRepository(db)

	// Expected query statement
	queryStatement := `SELECT * FROM public.users WHERE id=$1`

	// Filling mock database
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(int64(1), "Rizky")
	rows.AddRow(int64(2), "Ardi")
	rows.AddRow(int64(3), "Maulana")

	// Set expected outcome
	mock.ExpectQuery(regexp.QuoteMeta(queryStatement)).WithArgs(int64(1)).WillReturnRows(rows)

	// Excecute mock query
	_, err := Repository.GetUserById(context.Background(), int64(1))

	assert.NoError(t, err)
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func BenchmarkGetUsers(b *testing.B) {
	// creates sqlmock database connection and a mock to manage expectations
	db, mock := NewDBMock()

	// Setup Repository
	Repository := repository.NewUserRepository(db)

	// Expected query statement
	queryStatement := `SELECT * FROM public.users`

	// Filling mock database
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(int64(1), "Rizky")
	rows.AddRow(int64(2), "Ardi")
	rows.AddRow(int64(3), "Maulana")

	// Set expected outcome
	mock.ExpectQuery(regexp.QuoteMeta(queryStatement)).WillReturnRows(rows)

	// Excecute mock query
	_, err := Repository.GetUsers(context.Background())

	if err != nil {
		panic(err)
	}
}

func BenchmarkGetUserById(b *testing.B) {
	// creates sqlmock database connection and a mock to manage expectations
	db, mock := NewDBMock()

	// Setup Repository
	Repository := repository.NewUserRepository(db)

	// Expected query statement
	queryStatement := `SELECT * FROM public.users WHERE id=$1`

	// Filling mock database
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(int64(1), "Rizky")
	rows.AddRow(int64(2), "Ardi")
	rows.AddRow(int64(3), "Maulana")

	// Set expected outcome
	mock.ExpectQuery(regexp.QuoteMeta(queryStatement)).WithArgs(int64(1)).WillReturnRows(rows)

	// Excecute mock query
	user, err := Repository.GetUserById(context.Background(), int64(1))

	fmt.Println(user)
	if err != nil {
		panic(err)
	}
}
