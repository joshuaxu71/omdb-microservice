package common

import (
	"fmt"
	"os"
)

var OmdbAPI = "http://www.omdbapi.com/?apikey=%s&"

func InitOmdbAPI() {
	OmdbAPI = fmt.Sprintf(OmdbAPI, os.Getenv("OMDB_KEY"))
}

func IntegerAddress(value int) *int {
	integer := value
	return &integer
}
