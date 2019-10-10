package main

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type PropertyStore interface {
	Update(string, *Property, io.Reader) error
	GetProperty(string) (Property, error)
	GetPage(string) (io.ReadCloser, error)
}

type LocalPropertyStore struct {
	basePath string
}

func NewLocalPropertyStore(basePath string) LocalPropertyStore {
	os.MkdirAll(basePath, 0755)
	return LocalPropertyStore{basePath: basePath}
}

func (s *LocalPropertyStore) propertyPath(name string) string {
	return path.Join(s.basePath, name+".json")
}

func (s *LocalPropertyStore) pagePath(name string) string {
	return path.Join(s.basePath, name+".html")
}

func (s *LocalPropertyStore) updateProperty(name string, property *Property) error {
	now := time.Now()
	if property.FirstSeen.IsZero() {
		property.FirstSeen = now
	}
	property.LastSeen = now

	if f, err := os.Create(s.propertyPath(name)); err != nil {
		return err
	} else {
		defer f.Close()
		if b, err := json.MarshalIndent(property, "", "  "); err != nil {
			return err
		} else {
			f.Write(b)
		}
	}
	return nil
}

func (s *LocalPropertyStore) updatePage(name string, page io.Reader) error {
	if f, err := os.Create(s.pagePath(name)); err != nil {
		return err
	} else {
		defer f.Close()
		writer := bufio.NewWriter(f)
		if _, err := writer.ReadFrom(page); err != nil {
			return err
		}
		writer.Flush()
	}
	return nil
}

func (s *LocalPropertyStore) Update(name string, property *Property, page io.Reader) error {
	err := s.updateProperty(name, property)
	if err == nil {
		err = s.updatePage(name, page)
	}
	return err
}

func (s *LocalPropertyStore) GetProperty(name string) (Property, error) {
	var property Property
	if f, err := os.Open(s.propertyPath(name)); err != nil {
		return Property{}, err
	} else {
		defer f.Close()
		if data, err := ioutil.ReadAll(f); err != nil {
			return Property{}, err
		} else if err := json.Unmarshal(data, &property); err != nil {
			return Property{}, err
		}
	}
	return property, nil
}

func (s *LocalPropertyStore) GetPage(name string) (io.ReadCloser, error) {
	return os.Open(s.pagePath(name))
}
