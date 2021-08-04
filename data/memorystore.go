package data

import "sync"

var data = []Kitten{}

// MemoryStore is a simple in memory datastore that implements Store
type MemoryStore struct {
	nextid func() int
	sync.RWMutex
}

func NewMemoryStore() (*MemoryStore, error) {
	return &MemoryStore{nextid: GeneraId()}, nil
}

//Search returns a slice of Kitten which have a name matching the name in the parameters
func (m *MemoryStore) SearchName(name string) Kittens {
	defer m.RUnlock()
	var kittens Kittens
	m.RLock()
	for _, k := range data {
		if k.Name == name {
			kittens = append(kittens, k)
		}
	}
	return kittens
}

func (m *MemoryStore) DeleteAll() {
	defer m.Unlock()
	m.Lock()
	data = Kittens{}
}

// InsertKitten inserts a kitten into the datastore
func (m *MemoryStore) Insert(kitten Kitten) int {
	defer m.Unlock()
	m.Lock()
	i := m.nextid()
	kitten.Id = i
	data = append(data, kitten)
	return i
}

func (m *MemoryStore) DeleteName(name string) bool {
	defer m.Unlock()
	m.Lock()
	for i, k := range data {
		if k.Name == name {
			data = append(data[:i], data[i+1:]...)
			return true
		}
	}
	return false
}

func (m *MemoryStore) DeleteId(id int) bool {
	defer m.Unlock()
	m.Lock()
	for i, k := range data {
		if k.Id == id {
			data = append(data[:i], data[i+1:]...)
			return true
		}
	}
	return false
}

func (m *MemoryStore) List() Kittens {
	defer m.RUnlock()
	m.RLock()
	var result Kittens
	result = append(result, data...)
	return result
}

func (m *MemoryStore) SearchId(id int) *Kitten {
	defer m.RUnlock()
	m.RLock()
	for _, k := range data {
		if k.Id == id {
			kitten := k
			return &kitten
		}
	}
	return nil
}
