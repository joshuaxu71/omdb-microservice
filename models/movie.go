package models

type Movie struct {
	Title      string   `json:"Title"`
	Year       int      `json:"Year,string"`
	Rated      string   `json:"Rated,omitempty"`
	Released   string   `json:"Released,omitempty"`
	Runtime    string   `json:"Runtime,omitempty"`
	Genre      string   `json:"Genre,omitempty"`
	Director   string   `json:"Director,omitempty"`
	Writer     string   `json:"Writer,omitempty"`
	Actors     string   `json:"Actors,omitempty"`
	Plot       string   `json:"Plot,omitempty"`
	Language   string   `json:"Language,omitempty"`
	Country    string   `json:"Country,omitempty"`
	Awards     string   `json:"Awards,omitempty"`
	Poster     string   `json:"Poster"`
	Ratings    []Rating `json:"Ratings,omitempty"`
	Metascore  int      `json:"Metascore,string,omitempty"`
	ImdbRating float64  `json:"ImdbRating,string,omitempty"`
	ImdbVotes  int      `json:"ImdbVotes,string,omitempty"`
	ImdbId     string   `json:"ImdbId,omitempty"`
	Type       string   `json:"Type"`
	Dvd        string   `json:"Dvd,omitempty"`
	BoxOffice  string   `json:"BoxOffice,omitempty"`
	Production string   `json:"Production,omitempty"`
	Website    string   `json:"Website,omitempty"`
	Error      string   `json:"Error,omitempty"`
}

type SearchResult struct {
	Movies       []Movie `json:"Search"`
	TotalResults int     `json:"totalResults,string"`
}

type Rating struct {
	Source string `json:"Source"`
	Value  string `json:"Value"`
}
