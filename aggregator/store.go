package main

import (
	"fmt"
	"github.com/kirillApanasiuk/toll-calculator/types"
)

type MemoryStore struct {
	data map[int]float64
}

func (m *MemoryStore) Insert(distance types.Distance) error {
	m.data[distance.OBUID] += distance.Value
	return nil
}
func (m *MemoryStore) Get(id int) (float64, error) {
	dist, ok := m.data[id]
	if !ok {
		return 0.0, fmt.Errorf("could not find distance for obu id %d", id)
	}
	return dist, nil
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[int]float64),
	}
}
