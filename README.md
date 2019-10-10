# Ok Go Scan

```

type PropertyType int

const (
  None PropertyType = iota + 1
  Houses
  Flats
  Bungalows
  Land
  Commercial
  Other
)

Query {
  Postcode string
  Radius float
  MinBedrooms int
  MinPrice int 
  PropertyType PropertyType 
}


s := NewScanner(baseUrl)

q := NewQuery(func (q *Query) {
  q.Postcode = "W2 2SZ"
  q.PropertyType = Houses
  q.Radius = 1.0
  q.MinBedrooms = 2
})

properties, errs := s.Scan(q)

for {
  select {
    case property <- properties:
      fmt.Println(property)
    case err <- errs:
      fmt.Error(err)
      panic()
  }
}
```
