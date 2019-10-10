package main

import (
  "testing"
  "time"
  "bytes"
  "io/ioutil"
)

func TestUpdateProperty(t *testing.T) {
  // arrange
  store := NewLocalPropertyStore("./test-properties")
  property := Property {
    FirstSeen: time.Now(),
    LastSeen: time.Now(),
    Description: "dave dave",
    Address: "dave dave",
    Price: "123.45$$",
    Name: "foo-bar",
  }
  body := bytes.NewBuffer([]byte{'a'})

  // act
  err := store.Update("foo-bar", &property, body)

  // assert
  if err != nil {
    t.Error(err)
  }

  if property, err := store.GetProperty("foo-bar"); err != nil {
    t.Error(err)
  } else if property.Description != "dave dave" {
    t.Errorf("Invalid property object: %v", property)
  }

  if r, err := store.GetPage("foo-bar"); err != nil {
    t.Error(err)
  } else if body, err := ioutil.ReadAll(r); err != nil {
    t.Error(err)
  } else if string(body) != "a" {
    t.Error("invalid body")
  }

}
