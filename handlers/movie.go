package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"stock-bit/models"

	"github.com/gorilla/mux"
)

var OmdbAPI = "http://www.omdbapi.com/?apikey=%s&"

type MovieHandler struct {
	l *log.Logger
}

func NewMovieHandler(l *log.Logger) *MovieHandler {
	return &MovieHandler{l}
}

func InitOmdbAPI() {
	OmdbAPI = fmt.Sprintf(OmdbAPI, os.Getenv("OMDB_KEY"))
}

func (h *MovieHandler) GetMovieById(rw http.ResponseWriter, r *http.Request) {
	var uri string
	if mux.Vars(r)["id"] != "" {
		uri = fmt.Sprintf("%si=%s", OmdbAPI, mux.Vars(r)["id"])
	}
	if mux.Vars(r)["title"] != "" {
		uri = fmt.Sprintf("%st=%s", OmdbAPI, mux.Vars(r)["title"])
	}

	if uri == "" {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("id or title must have a value"))
		return
	}

	resp, err := http.Get(uri)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	defer resp.Body.Close()

	movie := models.Movie{}
	e := json.NewDecoder(resp.Body)
	e.Decode(&movie)
	if movie.Error != "" {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(movie.Error))
		return
	}

	marshalled, err := json.Marshal(movie)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	rw.Write(marshalled)
}

func (h *MovieHandler) GetMovies(rw http.ResponseWriter, r *http.Request) {
	var uri string

	if mux.Vars(r)["searchword"] == "" {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("searchword must have a value"))
		return
	}
	uri = fmt.Sprintf("%ss=%s", OmdbAPI, mux.Vars(r)["searchword"])

	if mux.Vars(r)["pagination"] != "" {
		uri = fmt.Sprintf("%s&page=%s", uri, mux.Vars(r)["pagination"])
	}

	resp, err := http.Get(uri)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	defer resp.Body.Close()

	searchResult := models.SearchResult{}
	e := json.NewDecoder(resp.Body)
	e.Decode(&searchResult)

	marshalled, err := json.Marshal(searchResult)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	rw.Write(marshalled)
}
