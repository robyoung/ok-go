package main

import (
	"bytes"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	// Instantiate default collector
	c := colly.NewCollector()
	c.IgnoreRobotsTxt = false
	c.Async = true
	c.UserAgent = "RobBot/0.1"
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       1 * time.Second,
	})

	store := NewLocalPropertyStore("./properties")

	c.OnHTML("body", func(e *colly.HTMLElement) {
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
		if err := store.Update(name, &property, bytes.NewReader(e.Response.Body)); err != nil {
			fmt.Println(err)
		}
	})

	// On every a element which has href attribute call callback
	c.OnHTML("a.propertyCard-priceLink[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %s\n", link)
		if propertyURL, err := url.Parse(e.Request.AbsoluteURL(link)); err == nil {
			c.Visit(propertyURL.String())
		}
	})

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

		fmt.Println("Scheduling index ", index, " of ", resultCount)
		c.Visit(newUrl.String())
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("")

	c.Wait()
}
