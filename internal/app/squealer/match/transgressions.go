package match

import (
	"fmt"
	"sync"
)

type Transgressions struct {
	sync.RWMutex
	internal map[string]*Transgression
	counter  int
}

func newTransgressions() *Transgressions {
	return &Transgressions{
		internal: make(map[string]*Transgression),
	}
}

func (t *Transgressions) Add(key string, transgression Transgression) {
	fmt.Println(transgression)

	t.Lock()
	if existing := t.internal[key]; existing == nil {
		t.internal[key] = &transgression
		t.counter += 1
		fmt.Printf("adding transgression %d\n", t.counter)
	}
	t.Unlock()
}

func (t *Transgressions) Exists(key string) bool {
	t.RLock()
	result := t.internal[key] != nil
	t.RUnlock()
	return result
}
