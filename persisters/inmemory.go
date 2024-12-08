package persisters

import "sync"

type InMemory struct {
	mu    sync.RWMutex
	store map[string]string
}

func (p *InMemory) Setup() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.store == nil {
		p.store = make(map[string]string)
	}
	return nil
}

func (p *InMemory) Set(key, value string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.store[key] = value
}

func (p *InMemory) Get(key string) string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.store[key]
}
