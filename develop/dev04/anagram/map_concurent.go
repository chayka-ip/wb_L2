package anagram

import "sync"

//Map is considered to be write-only during anagram processing
type mapConcurrent struct {
	data  map[string][]string
	mutex *sync.Mutex
}

func newMapConcurrent() *mapConcurrent {
	return &mapConcurrent{
		data:  map[string][]string{},
		mutex: &sync.Mutex{},
	}
}

func (m *mapConcurrent) AddValueSlice(key string, value []string) {
	m.mutex.Lock()
	m.data[key] = append(m.data[key], value...)
	m.mutex.Unlock()
}
