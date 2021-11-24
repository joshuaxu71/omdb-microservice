package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"stock-bit/common"
)

func GetMovieById(id string, title string) (*Movie, error) {
	var uri string
	if id != "" {
		uri = fmt.Sprintf("%si=%s", common.OmdbAPI, id)
	}
	if title != "" {
		uri = fmt.Sprintf("%st=%s", common.OmdbAPI, title)
	}

	if uri == "" {
		return nil, errors.New("id or title must have a value")
	}

	resp, err := http.Get(uri)
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
	return &movie, nil
}

func GetMovies(searchword string, pagination *int) (*SearchResult, error) {
	var uri string
	if searchword == "" {
		return nil, errors.New("searchword must have a value")
	}
	uri = fmt.Sprintf("%ss=%s", common.OmdbAPI, searchword)

	if pagination != nil {
		if *pagination < 1 {
			return nil, errors.New("pagination must be a larger number than 0")
		}
		uri = fmt.Sprintf("%s&page=%d", uri, *pagination)
	}

	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	searchResult := SearchResult{}
	e := json.NewDecoder(resp.Body)
	e.Decode(&searchResult)

	return &searchResult, nil
}
