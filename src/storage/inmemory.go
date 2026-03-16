package storage

import (
	"amber/src/encoder"
	"sync"
)

type inMemory struct {
	mu    sync.RWMutex
	store map[string]*Item
	list  []string
}

func NewInMemory() Service {
	return &inMemory{
		store: make(map[string]*Item),
		list:  []string{},
	}
}

func (m *inMemory) isUsed(id string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.store[id]
	return exists
}

func (m *inMemory) Save(url string) (string, error) {
	var id string
	l := 3

	for used := true; used; used = m.isUsed(id) {
		id = encoder.GetRandomId(l)
		l++
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	item := &Item{
		Id:     id,
		URL:    url,
		Visits: 0,
	}

	m.store[id] = item
	// Insert at the beginning to act like LPUSH
	m.list = append([]string{id}, m.list...)

	return id, nil
}

func (m *inMemory) Load(code string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, exists := m.store[code]
	if !exists {
		return "", ErrNoLink
	}

	item.Visits++
	return item.URL, nil
}

func (m *inMemory) LoadInfo(code string) (*Item, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.store[code]
	if !exists {
		return nil, ErrNoLink
	}

	// Return a copy to prevent external modification
	copiedItem := *item
	return &copiedItem, nil
}

func (m *inMemory) GetAll() []Item {
	m.mu.RLock()
	defer m.mu.RUnlock()

	links := make([]Item, 0, len(m.list))
	for _, id := range m.list {
		if item, exists := m.store[id]; exists {
			links = append(links, *item)
		}
	}

	return links
}

func (m *inMemory) Delete(code string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.store[code]; !exists {
		return ErrNoLink
	}

	delete(m.store, code)

	// Remove from list
	newList := []string{}
	for _, id := range m.list {
		if id != code {
			newList = append(newList, id)
		}
	}
	m.list = newList

	return nil
}

func (m *inMemory) Close() error {
	return nil
}
