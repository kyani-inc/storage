package local

import "sync"

type Local struct {
	store map[string][]byte
	mu    sync.RWMutex
}

func New() Local {
	return Local{store: map[string][]byte{}}
}

func (l *Local) Get(key string) []byte {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.store[key]
}

func (l *Local) Put(key string, data []byte) error {
	l.mu.Lock()
	l.store[key] = data
	l.mu.Unlock()
	return nil
}

func (l *Local) Delete(key string) {
	l.mu.Lock()
	delete(l.store, key)
	l.mu.Unlock()
}

func (l *Local) Flush() {
	l.mu.Lock()
	l.store = map[string][]byte{}
	l.mu.Unlock()
}
