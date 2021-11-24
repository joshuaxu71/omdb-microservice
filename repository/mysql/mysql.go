package mysql

import (
	"context"
	"database/sql"
	"time"

	repo "stock-bit/repository"

	_ "github.com/go-sql-driver/mysql"
)

// repository represent the repository model
type repository struct {
	db *sql.DB
}

// NewRepository will create a variable that represent the Repository struct
func NewRepository(dialect, dsn string, idleConn, maxConn int) (repo.Repository, error) {
	db, err := sql.Open(dialect, dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(idleConn)
	db.SetMaxOpenConns(maxConn)

	return &repository{db}, nil
}

// Close attaches the provider and close the connection
func (r *repository) Close() {
	r.db.Close()
}

// Find attaches the searchLog repository and find all data
func (r *repository) Find() ([]*repo.SearchLog, error) {
	searchLogs := make([]*repo.SearchLog, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, "SELECT type, url FROM search_logs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		searchLog := new(repo.SearchLog)
		err = rows.Scan(
			&searchLog.Type,
			&searchLog.URL,
		)

		if err != nil {
			return nil, err
		}
		searchLogs = append(searchLogs, searchLog)
	}

	return searchLogs, nil
}

// Create attaches the searchLog repository and creating the data
func (r *repository) Create(searchLog *repo.SearchLog) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "INSERT INTO search_logs (type, url) VALUES (?, ?)"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, searchLog.Type, searchLog.URL)
	return err
}
