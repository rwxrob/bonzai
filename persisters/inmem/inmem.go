/*
Package inmem, as the name implies, is an in-memory persistence layer with methods that are safe for concurrency. This is useful when simply wanting to work with something like a go map that is not safe for concurrency.
*/
package inmem

import "sync"

type Persister struct {
	mu    sync.RWMutex
	store map[string]string
}

// Setup creates the key-value map store safe for concurrency.
func (p *Persister) Setup() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.store == nil {
		p.store = make(map[string]string)
	}
	return nil
}

// Set stores a key-value pair in the key-value store safe for
// concurrency.
func (p *Persister) Set(key, value string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.store[key] = value
}

// Get retrieves the value associated with the given key from the
// key-value store safe for concurrency.
func (p *Persister) Get(key string) string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.store[key]
}
