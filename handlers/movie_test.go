package handlers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"reflect"
	"stock-bit/helpers"
	"stock-bit/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

var searchLogGetMovieByIdHTTP = &helpers.SearchLog{
	Transport: "HTTP",
	URL:       "http://test/movie",
}

var searchLogGetMoviesHTTP = &helpers.SearchLog{
	Transport: "HTTP",
	URL:       "http://test/movies",
}

var searchLogGetMovieByIdGRPC = &helpers.SearchLog{
	Transport: "gRPC",
	URL:       "http://test/movie",
}

var searchLogGetMoviesGRPC = &helpers.SearchLog{
	Transport: "gRPC",
	URL:       "http://test/movies",
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestGetMovieByIdHTTP(t *testing.T) {
	expectedResponse := models.Movie{
		Title:    "Doctor Strange",
		Year:     "2016",
		Rated:    "PG-13",
		Released: "04 Nov 2016",
		Runtime:  "115 min",
		Genre:    "Action, Adventure, Fantasy",
		Director: "Scott Derrickson",
		Writer:   "Jon Spaihts, Scott Derrickson, C. Robert Cargill",
		Actors:   "Benedict Cumberbatch, Chiwetel Ejiofor, Rachel McAdams",
		Plot:     "While on a journey of physical and spiritual healing, a brilliant neurosurgeon is drawn into the world of the mystic arts.",
		Language: "English",
		Country:  "United States",
		Awards:   "Nominated for 1 Oscar. 19 wins & 68 nominations total",
		Poster:   "https://m.media-amazon.com/images/M/MV5BNjgwNzAzNjk1Nl5BMl5BanBnXkFtZTgwMzQ2NjI1OTE@._V1_SX300.jpg",
		Ratings: []*models.Rating{
			{
				Source: "Internet Movie Database",
				Value:  "7.5/10",
			},
			{
				Source: "Rotten Tomatoes",
				Value:  "89%",
			},
			{
				Source: "Metacritic",
				Value:  "72/100",
			},
		},
		Metascore:  "72",
		ImdbRating: "7.5",
		ImdbVotes:  "654,213",
		ImdbId:     "tt1211837",
		Type:       "movie",
		Dvd:        "28 Feb 2017",
		BoxOffice:  "$232,641,920",
		Production: "N/A",
		Website:    "N/A",
	}

	srv := helpers.HttpMock("/movie", http.StatusOK, expectedResponse)
	defer srv.Close()

	ctx := context.WithValue(context.Background(), "transport", "HTTP")
	errs := make(chan error)
	db, mock := NewMock()
	query := `INSERT INTO search_logs(transport, url) VALUES (?, ?)`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(searchLogGetMovieByIdHTTP.Transport, searchLogGetMovieByIdHTTP.URL).
		WillReturnResult(sqlmock.NewResult(0, 1))
	go models.LogSearch(&ctx, errs, db, "http://test/movie")

	api := API{URL: srv.URL}
	response, err := api.GetMovieById()
	if err != nil {
		t.Error(err)
	}

	if err := <-errs; err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expectedResponse, *response) {
		t.Errorf("expected: %v\nactual: %v\n", expectedResponse, *response)
	}
}

func TestGetMovieByIdGRPC(t *testing.T) {
	expectedResponse := models.Movie{
		Title:    "Doctor Strange",
		Year:     "2016",
		Rated:    "PG-13",
		Released: "04 Nov 2016",
		Runtime:  "115 min",
		Genre:    "Action, Adventure, Fantasy",
		Director: "Scott Derrickson",
		Writer:   "Jon Spaihts, Scott Derrickson, C. Robert Cargill",
		Actors:   "Benedict Cumberbatch, Chiwetel Ejiofor, Rachel McAdams",
		Plot:     "While on a journey of physical and spiritual healing, a brilliant neurosurgeon is drawn into the world of the mystic arts.",
		Language: "English",
		Country:  "United States",
		Awards:   "Nominated for 1 Oscar. 19 wins & 68 nominations total",
		Poster:   "https://m.media-amazon.com/images/M/MV5BNjgwNzAzNjk1Nl5BMl5BanBnXkFtZTgwMzQ2NjI1OTE@._V1_SX300.jpg",
		Ratings: []*models.Rating{
			{
				Source: "Internet Movie Database",
				Value:  "7.5/10",
			},
			{
				Source: "Rotten Tomatoes",
				Value:  "89%",
			},
			{
				Source: "Metacritic",
				Value:  "72/100",
			},
		},
		Metascore:  "72",
		ImdbRating: "7.5",
		ImdbVotes:  "654,213",
		ImdbId:     "tt1211837",
		Type:       "movie",
		Dvd:        "28 Feb 2017",
		BoxOffice:  "$232,641,920",
		Production: "N/A",
		Website:    "N/A",
	}

	ctx := context.WithValue(context.Background(), "transport", "gRPC")
	errs := make(chan error)
	db, mock := NewMock()
	query := `INSERT INTO search_logs(transport, url) VALUES (?, ?)`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(searchLogGetMovieByIdGRPC.Transport, searchLogGetMovieByIdGRPC.URL).
		WillReturnResult(sqlmock.NewResult(0, 1))
	go models.LogSearch(&ctx, errs, db, "http://test/movie")

	mockMovieServiceClient := helpers.GrpcMock(t, "GetMovieById")
	response, err := mockMovieServiceClient.GetMovieById(
		context.Background(),
		&models.GetMovieByIdParams{Id: "tt1211837"},
	)
	if err != nil {
		t.Error(err)
	}

	if err := <-errs; err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expectedResponse, *response) {
		t.Errorf("expected: %v\nactual: %v\n", expectedResponse, *response)
	}
}

func TestGetMoviesHTTP(t *testing.T) {
	expectedResponse := models.SearchResult{
		Search: []*models.Movie{
			{
				Title:  "Doctor Strange",
				Year:   "2016",
				Poster: "https://m.media-amazon.com/images/M/MV5BNjgwNzAzNjk1Nl5BMl5BanBnXkFtZTgwMzQ2NjI1OTE@._V1_SX300.jpg",
				ImdbId: "tt1211837",
				Type:   "movie",
			},
			{
				Title:  "Strange Days",
				Year:   "1995",
				Poster: "https://m.media-amazon.com/images/M/MV5BODFkMTBmNjktMjM1Yy00MjY5LTliMGEtM2FhYjE2YjRmN2RkXkEyXkFqcGdeQXVyNzkwMjQ5NzM@._V1_SX300.jpg",
				ImdbId: "tt0114558",
				Type:   "movie",
			},
			{
				Title:  "Strange Wilderness",
				Year:   "2008",
				Poster: "https://m.media-amazon.com/images/M/MV5BMTI1MTI1OTUxNl5BMl5BanBnXkFtZTcwMDQzMjU1MQ@@._V1_SX300.jpg",
				ImdbId: "tt0489282",
				Type:   "movie",
			},
			{
				Title:  "The Adventures of Bob & Doug McKenzie: Strange Brew",
				Year:   "1983",
				Poster: "https://m.media-amazon.com/images/M/MV5BNGNmNmE1ZWQtNWYxNi00YjdjLWE2MTMtZjdiYzdmNjc2YzdjL2ltYWdlXkEyXkFqcGdeQXVyNzc5MjA3OA@@._V1_SX300.jpg",
				ImdbId: "tt0086373",
				Type:   "movie",
			},
			{
				Title:  "Jonathan Strange & Mr Norrell",
				Year:   "2015",
				Poster: "https://m.media-amazon.com/images/M/MV5BMTMzOGI1ZDAtYTk4NS00ODhlLTgxN2EtNDM1NDliNWM1M2UzL2ltYWdlXkEyXkFqcGdeQXVyNTAyODkwOQ@@._V1_SX300.jpg",
				ImdbId: "tt2548418",
				Type:   "series",
			},
			{
				Title:  "Love Is Strange",
				Year:   "2014",
				Poster: "https://m.media-amazon.com/images/M/MV5BMTk5MTkxOTI1N15BMl5BanBnXkFtZTgwNzAwNDA4MTE@._V1_SX300.jpg",
				ImdbId: "tt2639344",
				Type:   "movie",
			},
			{
				Title:  "Strange Magic",
				Year:   "2015",
				Poster: "https://m.media-amazon.com/images/M/MV5BMjA0NjU3MTU5OF5BMl5BanBnXkFtZTgwMTYyMDQ3MzE@._V1_SX300.jpg",
				ImdbId: "tt4191054",
				Type:   "movie",
			},
			{
				Title:  "Life Is Strange",
				Year:   "2015",
				Poster: "https://m.media-amazon.com/images/M/MV5BNzc4MjI0NDgtMWM2ZC00NmMxLThmYjMtYmM3M2ZlYzhlNzVmXkEyXkFqcGdeQXVyMjYwNDA2MDE@._V1_SX300.jpg",
				ImdbId: "tt4375662",
				Type:   "game",
			},
			{
				Title:  "When You're Strange",
				Year:   "2009",
				Poster: "https://m.media-amazon.com/images/M/MV5BMTg1NjQ3OTQwOF5BMl5BanBnXkFtZTcwMDE5NTgyMw@@._V1_SX300.jpg",
				ImdbId: "tt1333667",
				Type:   "movie",
			},
			{
				Title:  "The Strange Love of Martha Ivers",
				Year:   "1946",
				Poster: "https://m.media-amazon.com/images/M/MV5BZmJjMDIxMGUtZWM1NS00ZTg5LWI2NWQtNjUxZjc5MjI5OTc3XkEyXkFqcGdeQXVyNjc0MzMzNjA@._V1_SX300.jpg",
				ImdbId: "tt0038988",
				Type:   "movie",
			},
		},
		TotalResults: "989",
	}

	srv := helpers.HttpMock("/movies", http.StatusOK, expectedResponse)
	defer srv.Close()

	ctx := context.WithValue(context.Background(), "transport", "HTTP")
	errs := make(chan error)
	db, mock := NewMock()
	query := `INSERT INTO search_logs(transport, url) VALUES (?, ?)`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(searchLogGetMoviesHTTP.Transport, searchLogGetMoviesHTTP.URL).
		WillReturnResult(sqlmock.NewResult(0, 1))
	go models.LogSearch(&ctx, errs, db, "http://test/movies")

	api := API{URL: srv.URL}
	response, err := api.GetMovies()
	if err != nil {
		t.Error(err)
	}

	if err := <-errs; err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expectedResponse, *response) {
		t.Errorf("expected: %v\nactual: %v\n", expectedResponse, *response)
	}
}

func TestGetMoviesGRPC(t *testing.T) {
	expectedResponse := models.SearchResult{
		Search: []*models.Movie{
			{
				Title:  "Doctor Strange",
				Year:   "2016",
				Poster: "https://m.media-amazon.com/images/M/MV5BNjgwNzAzNjk1Nl5BMl5BanBnXkFtZTgwMzQ2NjI1OTE@._V1_SX300.jpg",
				ImdbId: "tt1211837",
				Type:   "movie",
			},
			{
				Title:  "Strange Days",
				Year:   "1995",
				Poster: "https://m.media-amazon.com/images/M/MV5BODFkMTBmNjktMjM1Yy00MjY5LTliMGEtM2FhYjE2YjRmN2RkXkEyXkFqcGdeQXVyNzkwMjQ5NzM@._V1_SX300.jpg",
				ImdbId: "tt0114558",
				Type:   "movie",
			},
			{
				Title:  "Strange Wilderness",
				Year:   "2008",
				Poster: "https://m.media-amazon.com/images/M/MV5BMTI1MTI1OTUxNl5BMl5BanBnXkFtZTcwMDQzMjU1MQ@@._V1_SX300.jpg",
				ImdbId: "tt0489282",
				Type:   "movie",
			},
			{
				Title:  "The Adventures of Bob & Doug McKenzie: Strange Brew",
				Year:   "1983",
				Poster: "https://m.media-amazon.com/images/M/MV5BNGNmNmE1ZWQtNWYxNi00YjdjLWE2MTMtZjdiYzdmNjc2YzdjL2ltYWdlXkEyXkFqcGdeQXVyNzc5MjA3OA@@._V1_SX300.jpg",
				ImdbId: "tt0086373",
				Type:   "movie",
			},
			{
				Title:  "Jonathan Strange & Mr Norrell",
				Year:   "2015",
				Poster: "https://m.media-amazon.com/images/M/MV5BMTMzOGI1ZDAtYTk4NS00ODhlLTgxN2EtNDM1NDliNWM1M2UzL2ltYWdlXkEyXkFqcGdeQXVyNTAyODkwOQ@@._V1_SX300.jpg",
				ImdbId: "tt2548418",
				Type:   "series",
			},
			{
				Title:  "Love Is Strange",
				Year:   "2014",
				Poster: "https://m.media-amazon.com/images/M/MV5BMTk5MTkxOTI1N15BMl5BanBnXkFtZTgwNzAwNDA4MTE@._V1_SX300.jpg",
				ImdbId: "tt2639344",
				Type:   "movie",
			},
			{
				Title:  "Strange Magic",
				Year:   "2015",
				Poster: "https://m.media-amazon.com/images/M/MV5BMjA0NjU3MTU5OF5BMl5BanBnXkFtZTgwMTYyMDQ3MzE@._V1_SX300.jpg",
				ImdbId: "tt4191054",
				Type:   "movie",
			},
			{
				Title:  "Life Is Strange",
				Year:   "2015",
				Poster: "https://m.media-amazon.com/images/M/MV5BNzc4MjI0NDgtMWM2ZC00NmMxLThmYjMtYmM3M2ZlYzhlNzVmXkEyXkFqcGdeQXVyMjYwNDA2MDE@._V1_SX300.jpg",
				ImdbId: "tt4375662",
				Type:   "game",
			},
			{
				Title:  "When You're Strange",
				Year:   "2009",
				Poster: "https://m.media-amazon.com/images/M/MV5BMTg1NjQ3OTQwOF5BMl5BanBnXkFtZTcwMDE5NTgyMw@@._V1_SX300.jpg",
				ImdbId: "tt1333667",
				Type:   "movie",
			},
			{
				Title:  "The Strange Love of Martha Ivers",
				Year:   "1946",
				Poster: "https://m.media-amazon.com/images/M/MV5BZmJjMDIxMGUtZWM1NS00ZTg5LWI2NWQtNjUxZjc5MjI5OTc3XkEyXkFqcGdeQXVyNjc0MzMzNjA@._V1_SX300.jpg",
				ImdbId: "tt0038988",
				Type:   "movie",
			},
		},
		TotalResults: "989",
	}

	ctx := context.WithValue(context.Background(), "transport", "gRPC")
	errs := make(chan error)
	db, mock := NewMock()
	query := `INSERT INTO search_logs(transport, url) VALUES (?, ?)`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(searchLogGetMoviesGRPC.Transport, searchLogGetMoviesGRPC.URL).
		WillReturnResult(sqlmock.NewResult(0, 1))
	go models.LogSearch(&ctx, errs, db, "http://test/movies")

	mockMovieServiceClient := helpers.GrpcMock(t, "GetMovies")
	response, err := mockMovieServiceClient.GetMovies(
		context.Background(),
		&models.GetMoviesParams{
			Searchword: "strange",
			Pagination: "1",
		},
	)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expectedResponse, *response) {
		t.Errorf("expected: %v\nactual: %v\n", expectedResponse, *response)
	}
}
