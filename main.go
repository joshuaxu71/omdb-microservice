package main

import (
	"log"
	"net/http"
	"os"
	"stock-bit/handlers"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// loading the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// create the logger
	l := log.New(os.Stdout, "omdb-api: ", log.LstdFlags)

	// create the router
	r := mux.NewRouter()

	// add routes to the router
	handlers.InitOmdbAPI()
	handler := handlers.NewMovieHandler(l)

	getRouter := r.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/movie", handler.GetMovieById).
		Queries(
			"id", "{id}",
		)
	getRouter.HandleFunc("/movie", handler.GetMovieById).
		Queries(
			"title", "{title}",
		)
	getRouter.HandleFunc("/movies", handler.GetMovies).
		Queries(
			"searchword", "{searchword}",
			"pagination", "{pagination}",
		)
	getRouter.HandleFunc("/movies", handler.GetMovies).
		Queries(
			"searchword", "{searchword}",
		)

	// create a new server
	s := http.Server{
		Addr:         os.Getenv("PORT"), // configure the bind address
		Handler:      r,                 // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	l.Printf("Starting server on port %s", os.Getenv("PORT"))

	err = s.ListenAndServe()
	if err != nil {
		l.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
