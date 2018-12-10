package batcher

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	COUNT = 5
)

func invoke(b *Batcher, ch chan int) {
	b.Invoke()
	ch <- 0
}

func TestSingleSyncInvocation(t *testing.T) {
	count := 0
	b := NewBatcher(func() {
		count++
	})

	b.Invoke()

	assert.Equal(t, 1, count)
}

func TestMultipleSyncInvocations(t *testing.T) {
	count := 0
	b := NewBatcher(func() {
		count++
	})

	for i := 0; i < COUNT; i++ {
		b.Invoke()
	}

	assert.Equal(t, COUNT, count)
}

func TestSingleAsyncInvocation(t *testing.T) {
	count := 0
	b := NewBatcher(func() {
		count++
	})
	gate := make(chan int, 1)

	go func() {
		b.Invoke()
		gate <- 0
	}()
	<-gate

	assert.Equal(t, 1, count)
}

func TestMultipleAsyncInvocations(t *testing.T) {
	gate1 := make(chan int, COUNT)
	gate2 := make(chan int, COUNT)
	count := 0
	b := NewBatcher(func() {
		<-gate1
		count++
	})

	for i := 0; i < COUNT; i++ {
		go func() {
			b.Invoke()
			gate2 <- 0
		}()
	}

	time.Sleep(100 * time.Millisecond)
	for i := 0; i < COUNT; i++ {
		gate1 <- 0
	}
	for i := 0; i < COUNT; i++ {
		<-gate2
	}

	assert.Equal(t, 1, count)
}

func TestSequenceOfInvocations(t *testing.T) {
	gate1 := make(chan int, COUNT)
	gate2 := make(chan int, COUNT)
	count := 0
	b := NewBatcher(func() {
		<-gate1
		count++
	})

	for i := 0; i < COUNT; i++ {
		go func() {
			b.Invoke()
			gate2 <- 0
		}()
	}
	time.Sleep(100 * time.Millisecond)
	gate1 <- 0
	for i := 0; i < COUNT; i++ {
		<-gate2
	}

	for i := 0; i < COUNT; i++ {
		go func() {
			b.Invoke()
			gate2 <- 0
		}()
	}
	time.Sleep(100 * time.Millisecond)
	gate1 <- 0
	for i := 0; i < COUNT; i++ {
		<-gate2
	}

	assert.Equal(t, 2, count)
}

func TestPanicInInvocation(t *testing.T) {
	testErr := struct{ tag string }{tag: "Test"}
	gate := make(chan interface{})
	b := NewBatcher(func() {
		panic(testErr)
	})

	go func() {
		err := b.Invoke()
		gate <- err
	}()
	err := <-gate

	assert.True(t, err == testErr)
}

func TestPanicInSeveralIncocations(t *testing.T) {
	testErr := struct{ tag string }{tag: "Test"}
	gate1 := make(chan int, COUNT)
	gate2 := make(chan interface{}, COUNT)
	b := NewBatcher(func() {
		panic(testErr)
	})

	for i := 0; i < COUNT; i++ {
		go func() {
			err := b.Invoke()
			gate2 <- err
		}()
	}

	time.Sleep(100 * time.Millisecond)
	for i := 0; i < COUNT; i++ {
		gate1 <- 0
	}
	var errs []interface{}
	for i := 0; i < COUNT; i++ {
		err := <-gate2
		errs = append(errs, err)
	}

	assert.Equal(t, []interface{}{testErr, testErr, testErr, testErr, testErr}, errs)
}
