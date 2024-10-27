package main

import "sync"

// Lehmer is a xorshift+ generator
type Lehmer struct {
	last uint64
	lock sync.Mutex
}

const lehmerA uint64 = 279470273
const lehmerB uint64 = 4294967291

// NewLehmer creates a Lehmer generator
func NewLehmer(seed int64) *Lehmer {
	gen := &Lehmer{}
	gen.Seed(seed)
	return gen
}

// Seed seeds the generator
func (gen *Lehmer) Seed(seed int64) {
	gen.lock.Lock()
	defer gen.lock.Unlock()

	gen.last = uint64(seed)
	if gen.last == 0 {
		gen.last = 1
	}
}

// Int63 returns a pseudo-random int64
func (gen *Lehmer) Int63() int64 {
	gen.lock.Lock()
	defer gen.lock.Unlock()

	gen.last = ((gen.last * lehmerA) % lehmerB) + 1
	return int64(gen.last - 1)
}
