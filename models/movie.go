package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"stock-bit/common"
)

// a function that will be ran by a routine to log search calls into MySQL db
func LogSearch(c *context.Context, errs chan error, db *sql.DB, url string) {
	ctx := *c
	query := "INSERT INTO search_logs(transport, url) VALUES (?, ?)"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		errs <- fmt.Errorf("error %s when preparing SQL statement", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, ctx.Value("transport"), url)
	if err != nil {
		errs <- fmt.Errorf("error %s when inserting row into search_logs table", err)
	}

	errs <- nil
}

func GetMovieById(ctx *context.Context, id string, title string) (*Movie, error) {
	var url string
	if id != "" {
		url = fmt.Sprintf("%si=%s", common.OmdbAPI, id)
	}
	if title != "" {
		url = fmt.Sprintf("%st=%s", common.OmdbAPI, title)
	}

	if url == "" {
		return nil, errors.New("id or title must have a value")
	}

	errs := make(chan error)
	db := common.GetDb()
	go LogSearch(ctx, errs, db, url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	movie := Movie{}
	e := json.NewDecoder(resp.Body)
	e.Decode(&movie)
	if movie.Error != "" {
		return nil, errors.New(movie.Error)
	}

	if err := <-errs; err != nil {
		return nil, err
	}

	return &movie, nil
}

func GetMovies(ctx *context.Context, searchword string, pagination *int) (*SearchResult, error) {
	var url string
	if searchword == "" {
		return nil, errors.New("searchword must have a value")
	}
	url = fmt.Sprintf("%ss=%s", common.OmdbAPI, searchword)

	if pagination != nil {
		if *pagination < 1 {
			return nil, errors.New("pagination must be a larger number than 0")
		}
		url = fmt.Sprintf("%s&page=%d", url, *pagination)
	}

	errs := make(chan error)
	db := common.GetDb()
	go LogSearch(ctx, errs, db, url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	searchResult := SearchResult{}
	e := json.NewDecoder(resp.Body)
	e.Decode(&searchResult)

	if err := <-errs; err != nil {
		return nil, err
	}

	return &searchResult, nil
}
