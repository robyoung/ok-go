package main

import (
  "testing"
)

func TestQuerySetup(t *testing.T) {
  q1 := NewQuery(func (q *Query) {
    q.PropertyType = Houses
  })
  q2 := NewQuery(func (q *Query) {})

  if q1.PropertyType != Houses {
    t.Errorf("Expected property type of houses got %s", q1.PropertyType.String())
  }
  if q2.PropertyType != 0 {
    t.Errorf("Expected property type of none got %s", q2.PropertyType.String())
  }
}
