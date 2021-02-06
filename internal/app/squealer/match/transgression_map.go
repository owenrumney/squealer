package match

import (
	"sync"
)

type transgressionMap struct {
	sync.RWMutex
	internal map[string]*Transgression
	counter  int
}

func newTransgressions() *transgressionMap {
	return &transgressionMap{
		internal: make(map[string]*Transgression),
	}
}

func (t *transgressionMap) add(key string, transgression Transgression) {
	t.Lock()
	if existing := t.internal[key]; existing == nil {
		t.internal[key] = &transgression
		t.counter += 1
	}
	t.Unlock()
}

func (t *transgressionMap) exists(key string) bool {
	t.RLock()
	result := t.internal[key] != nil
	t.RUnlock()
	return result
}

func (t *transgressionMap) count() int {
	return t.counter
}
