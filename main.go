package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	query := NewQuery(func(q *Query) {
		q.Postcode = "W2 2SZ"
		q.Radius = 1.0
	})
	scanner := NewScanner(os.Args[1])
	properties, errs := scanner.Scan(query)
	for {
		select {
		case property := <-properties:
			if property.IsEmpty() {
				fmt.Println("empty property")
			} else {
				fmt.Println(property.OneLine())
			}
		case err := <-errs:
			if err == nil {
				return
			} else {
				log.Fatal(err)
			}
		}
	}
}
