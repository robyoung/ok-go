package main

import (
  "time"
)

type Property struct {
  FirstSeen time.Time `json:"firstSeen"`
  LastSeen time.Time `json:"lastSeen"`
  Description string `json:"description"`
  Address string `json:"address"`
  Price string `json:"price"`
  Name string `json:-`
}
