package repository

type Repository interface {
	Close()
	Find() ([]*SearchLog, error)
	Create(searchLog *SearchLog) error
}

type SearchLog struct {
	Type string
	URL  string
}
