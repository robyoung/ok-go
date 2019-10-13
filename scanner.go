package main

import (
	"net/url"
  "log"
  "fmt"
  "time"
  "strconv"
  "strings"
  "path"

	"github.com/gocolly/colly"
)

type Scanner struct {
	baseURL string
}

func NewScanner(baseURL string) Scanner {
	return Scanner{baseURL: baseURL}
}

func (s *Scanner) Scan(query Query) (<-chan Property, <-chan error) {
	properties := make(chan Property, 1)
	errs := make(chan error, 1)

	c := colly.NewCollector()
	c.IgnoreRobotsTxt = false
	c.Async = true
	c.UserAgent = "RobBot/0.1"
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       1 * time.Second,
	})

  // First stage search
  c.OnHTML("form#initialSearch", func(e *colly.HTMLElement) {
    fields := s.getInputFields(e)

    // set searchLocation
    fields["searchLocation"] = query.Postcode

    // set button action
    fields["buy"] = "For sale"

    action, err := s.makeQueryURL(e.Attr("action"), fields)
    if err != nil {
      log.Fatal("Invalid first stage action")
    }

    fmt.Printf("Search Query 1: %s\n", action.String())

    c.Visit(action.String())
  })

  // Second stage search
  c.OnHTML("form#propertySearchCriteria", func(e *colly.HTMLElement) {

    fields := s.getInputFields(e)

    fields["radius"] = fmt.Sprintf("%.1f", query.Radius)
    fields["minBedrooms"] = s.formatInteger(query.MinBedrooms)
    fields["maxBedrooms"] = ""
    fields["minPrice"] = s.formatInteger(query.MinPrice)
    fields["maxPrice"] = ""
    fields["displayPropertyType"] = query.PropertyType.String()
    fields["maxDaysSinceAdded"] = ""

    action, err := s.makeQueryURL(e.Attr("action"), fields)
    if err != nil {
      log.Fatal("Invalid second stage action")
    }

    fmt.Printf("Search Query 2: %s\n", action.String())

    c.Visit(action.String())
  })

  // Pagination
	c.OnHTML("span.searchHeader-resultCount", func(e *colly.HTMLElement) {
		var err error
		var resultCount int
		var index int

		if resultCount, err = strconv.Atoi(e.Text); err != nil {
			return
		}

		if queryIndex := e.Request.URL.Query().Get("index"); queryIndex != "" {
			if index, err = strconv.Atoi(queryIndex); err != nil {
				return
			}
		}

		if index < resultCount {
			index += 24
		}

		newUrl, _ := url.Parse(e.Request.URL.String())
		query := newUrl.Query()
		query.Set("index", strconv.Itoa(index))
		newUrl.RawQuery = query.Encode()

    fmt.Println("Next page: ", newUrl.String())
		c.Visit(newUrl.String())
	})

  // Property URL scheduling
	c.OnHTML("a.propertyCard-priceLink[href]", func(e *colly.HTMLElement) {
		mainURL := e.Request.AbsoluteURL(e.Attr("href"))
		// Print link
		if propertyURL, err := url.Parse(mainURL); err == nil {
			c.Visit(propertyURL.String())
		}
	})

  // Property page
	c.OnHTML("body", func(e *colly.HTMLElement) {
    fmt.Println("Visiting: ", e.Request.URL.String())
		if !strings.Contains(e.Request.URL.Path, "property-for-sale/property") {
			return
		}
		name := path.Base(e.Request.URL.Path)
		ext := path.Ext(name)
		if ext != "" {
			name = name[:len(name)-len(ext)]
		}
		propElem := e.DOM.Find("div.property-header-bedroom-and-price")
		property := Property{}
		property.Description = propElem.Find("h1").Text()
		property.Address = strings.TrimSpace(e.DOM.Find("address").Text())
		property.Price = strings.TrimSpace(e.DOM.Find("#propertyHeaderPrice strong").Text())
		property.Name = name
    property.Page = e.Response.Body

    properties <- property
	})

  go func() {
    c.Visit(s.baseURL)
    c.Wait()
    close(properties)
    close(errs)
  }()


	return properties, errs
}

func (s *Scanner) parseAction(action string) (target *url.URL, err error) {
  target, err = url.Parse(s.baseURL)
  if err != nil {
    return
  }
  target, err = target.Parse(action)
  return
}

func (s *Scanner) makeQueryURL(action string, fields map[string]string) (target *url.URL, err error) {
  target, err = s.parseAction(action)
  if err != nil {
    return
  }

  // Make the query
  targetQuery := url.Values{}
  for k, v := range fields {
    targetQuery.Add(k, v)
  }
  target.RawQuery = targetQuery.Encode()
  return
}

func (s *Scanner) formatInteger(i int) string {
  if i < 1 {
    return ""
  } else {
    return fmt.Sprintf("%d", i)
  }
}

func (s *Scanner) getInputFields(e *colly.HTMLElement) map[string]string {
  fields := make(map[string]string)
  e.ForEach("input", func(i int, e1 *colly.HTMLElement) {
    fields[e1.Attr("name")] = e1.Attr("value")
  })

  return fields
}


func newCollector(baseURL string) *colly.Collector {
  return colly.NewCollector()
}



type PropertyType int

const (
  Any PropertyType = iota + 1
	Houses
	Flats
	Bungalows
	Land
	Commercial
	Other
)

func (p PropertyType) String() string {
	switch p {
  case Any:
    return ""
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
		return ""
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
