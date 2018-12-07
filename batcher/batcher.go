// Package batcher contains batcher type that replaces several similar long running
// invocations with a single one.
package batcher

import (
	"sync"
)

// Batcher replaces several long running function invocations with a single one.
type Batcher struct {
	wg     sync.WaitGroup
	mux    sync.Mutex
	locker int
	action func()
}

// NewBatcher creates an instance of Batcher.
func NewBatcher(action func()) *Batcher {
	return &Batcher{action: action}
}

func (b *Batcher) lock() bool {
	b.mux.Lock()
	defer b.mux.Unlock()
	if b.locker < 0 {
		panic("negative locker")
	}
	if b.locker > 0 {
		b.locker++
		return false
	}
	b.locker = 1
	b.wg.Add(1)
	return true
}

func (b *Batcher) unlock() {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.locker--
	if b.locker < 0 {
		panic("negative locker")
	}
}

func (b *Batcher) call() {
	defer b.wg.Done()
	b.action()
}

// Invoke calls batcher action.
func (b *Batcher) Invoke() {
	if b.lock() {
		b.call()
	}
	defer b.unlock()
	b.wg.Wait()
}
