package mysql

import (
	"database/sql"
	"log"
	"testing"

	r "stock-bit/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var searchLog = &r.SearchLog{
	Type: uuid.New().String(),
	URL:  "Momo",
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}
func TestFind(t *testing.T) {
	db, mock := NewMock()
	repo := &repository{db}
	defer func() {
		repo.Close()
	}()

	query := "SELECT type, url FROM search_logs"

	rows := sqlmock.NewRows([]string{"type", "url"}).
		AddRow(searchLog.Type, searchLog.URL)

	mock.ExpectQuery(query).WillReturnRows(rows)

	searchLogs, err := repo.Find()
	assert.NotEmpty(t, searchLogs)
	assert.NoError(t, err)
	assert.Len(t, searchLogs, 1)
}

func TestCreate(t *testing.T) {
	db, mock := NewMock()
	repo := &repository{db}
	defer func() {
		repo.Close()
	}()

	query := "INSERT INTO search_logs \\(type, url\\) VALUES \\(\\?, \\?\\)"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(searchLog.Type, searchLog.URL).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Create(searchLog)
	assert.NoError(t, err)
}

func TestCreateError(t *testing.T) {
	db, mock := NewMock()
	repo := &repository{db}
	defer func() {
		repo.Close()
	}()

	query := "INSERT INTO search_logs \\(type, url\\) VALUES \\(\\?, \\?\\)"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(searchLog.Type, searchLog.URL).WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Create(searchLog)
	assert.Error(t, err)
}
