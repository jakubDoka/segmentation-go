package nodes

import (
	"errors"
)

// Item is content of storage
type Item int

// Item enum
const (
	TIN Item = iota
	IRON
	STEEL
	CARBON
)

// Items stores all Items from enum
var Items []Item = []Item{
	TIN,
	IRON,
	STEEL,
	CARBON,
}

// Storage stores items and manages item requests
type Storage struct {
	Stored, Requested map[Item]int
}

// NewStorage returns new storage ready to use
func NewStorage() *Storage {
	return &Storage{
		Stored:    map[Item]int{},
		Requested: map[Item]int{},
	}
}

// Add adds item to storage
func (s *Storage) Add(item Item, amount int) {
	s.Stored[item] += amount
	if s.Requested[item] < 0 {
		panic(errors.New("negative request value"))
	}
}

// Merge merges two storages
func (s *Storage) Merge(o *Storage) {
	for k, v := range o.Stored {
		s.Add(k, v)
		o.Stored[k] = 0
	}
}

// IsEmpty returns weather storage is not empty
func (s *Storage) IsEmpty() bool {
	for _, v := range s.Stored {
		if v != 0 {
			return false
		}
	}
	return true
}

// Use uses items in storage
func (s *Storage) Use(item Item, amount int) {
	s.Stored[item] -= amount
	s.Requested[item] += amount
}
