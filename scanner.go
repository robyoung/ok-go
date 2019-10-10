package main

type Scanner struct {
	baseURL string
}

func NewScanner(baseURL string) Scanner {
	return Scanner{baseURL: baseURL}
}

func (s *Scanner) Scan(query Query) (<-chan Property, <-chan error) {
	properties := make(chan Property, 1)
	errs := make(chan error, 1)

	return properties, errs
}

type PropertyType int

const (
	Houses PropertyType = iota + 1
	Flats
	Bungalows
	Land
	Commercial
	Other
)

func (p PropertyType) String() string {
	switch p {
	case Houses:
		return "houses"
	case Flats:
		return "flats"
	case Bungalows:
		return "bungalows"
	case Land:
		return "land"
	case Commercial:
		return "commercial"
	case Other:
		return "other"
	default:
		return "none"
	}
}

type Query struct {
	Postcode     string
	Radius       float32
	MinBedrooms  int
	MinPrice     int
	PropertyType PropertyType
}

func NewQuery(setup func(*Query)) Query {
	query := Query{}
	setup(&query)
	return query
}
