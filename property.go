package main

import (
	"fmt"
	"strings"
	"time"
)

type Property struct {
	FirstSeen   time.Time `json:"firstSeen"`
	LastSeen    time.Time `json:"lastSeen"`
	Description string    `json:"description"`
	Address     string    `json:"address"`
	Price       string    `json:"price"`
	Page        []byte    `json:"-"`
	Name        string    `json:"-"`
}

func (p *Property) IsEmpty() bool {
	if p.FirstSeen != (time.Time{}) {
		return false
	}
	if p.LastSeen != (time.Time{}) {
		return false
	}
	if p.Description != "" {
		return false
	}
	if p.Address != "" {
		return false
	}
	if p.Price != "" {
		return false
	}
	return true
}

func (p *Property) OneLine() string {
	return fmt.Sprintf("%s (%s)", strings.ReplaceAll(p.Description, "\n", ""), p.Price)
}
