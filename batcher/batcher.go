// Package batcher contains batcher type that replaces several similar long running
// invocations with a single one.
package batcher

import (
	"sync"
)

/*
Batcher replaces several long running function invocations with a single one.

	func DoSomeTask() {
		... // time.Sleep(5 * time.Second)
	}

	for i := 0; i < COUNT; i++ {
		go func() {
			DoSomeTask()
		}()
		...
	}

Let there be some function that is called from several threads.
The function is either long running or resource consuming in any other way.
It would be better if all threads reuse single function invocation.

Batcher does it.

	batcher := NewBatcher(DoSomeTask)

	for i := 0; i < COUNT; i++ {
		go func() {
			batcher.Invoke()
		}()
		...
	}

*/
type Batcher struct {
	wg     *sync.WaitGroup
	mux    *sync.Mutex
	locker int
	err    interface{}
	action func()
}

// NewBatcher creates an instance of Batcher.
func NewBatcher(action func()) *Batcher {
	if action == nil {
		panic("nil action")
	}
	return &Batcher{
		wg:     &sync.WaitGroup{},
		mux:    &sync.Mutex{},
		action: action,
	}
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
	defer b.capturePanic()
	b.action()
}

func (b *Batcher) capturePanic() {
	b.err = recover()
}

// Invoke executes batched action.
func (b *Batcher) Invoke() interface{} {
	if b.lock() {
		b.call()
	}
	defer b.unlock()
	b.wg.Wait()
	return b.err
}
