package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"stock-bit/common"
	"stock-bit/helpers"
	"stock-bit/models"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type MovieHandler struct {
	l *log.Logger
}

func NewMovieHandler(l *log.Logger) *MovieHandler {
	return &MovieHandler{l}
}

func (h *MovieHandler) GetMovieById(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, "transport", "HTTP")
	movie, err := models.GetMovieById(&ctx, mux.Vars(r)["id"], mux.Vars(r)["title"])
	if err != nil {
		rw.Write([]byte(err.Error()))
		return
	} else if movie == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(err.Error()))
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
	if mux.Vars(r)["pagination"] == "" {
		mux.Vars(r)["pagination"] = "1"
	}

	pagination, err := strconv.Atoi(mux.Vars(r)["pagination"])
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("pagination has to be an integer"))
		return
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, "transport", "HTTP")
	searchResult, err := models.GetMovies(&ctx, mux.Vars(r)["searchword"], common.IntegerAddress(pagination))
	if err != nil {
		rw.Write([]byte(err.Error()))
		return
	}

	marshalled, err := json.Marshal(searchResult)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	rw.Write(marshalled)
}

type Server struct {
}

func (s *Server) GetMovieById(ctx context.Context, in *models.GetMovieByIdParams) (*models.Movie, error) {
	ctx = context.WithValue(ctx, "transport", "gRPC")
	movie, err := models.GetMovieById(&ctx, in.Id, in.Title)
	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (s *Server) GetMovies(ctx context.Context, in *models.GetMoviesParams) (*models.SearchResult, error) {
	if in.Pagination == "" {
		in.Pagination = "1"
	}

	pagination, err := strconv.Atoi(in.Pagination)
	if err != nil {
		return nil, errors.New("pagination has to be an integer")
	}

	ctx = context.WithValue(ctx, "transport", "gRPC")
	searchResult, err := models.GetMovies(&ctx, in.Searchword, common.IntegerAddress(int(pagination)))
	if err != nil {
		return nil, err
	}

	return searchResult, nil
}

type API struct {
	URL string
}

func (api *API) GetMovieById() (*models.Movie, error) {
	movie := &models.Movie{}

	to := time.Duration(10)
	opt := &helpers.HttpOptions{
		Ctx:    context.Background(),
		Url:    api.URL + "/movie",
		TO:     &to,
		Method: http.MethodGet,
	}

	_, err := helpers.DoRequest(opt, movie)
	return movie, err
}

func (api *API) GetMovies() (*models.SearchResult, error) {
	searchResult := &models.SearchResult{}

	to := time.Duration(10)
	opt := &helpers.HttpOptions{
		Ctx:    context.Background(),
		Url:    api.URL + "/movies",
		TO:     &to,
		Method: http.MethodGet,
	}

	_, err := helpers.DoRequest(opt, searchResult)
	return searchResult, err
}
