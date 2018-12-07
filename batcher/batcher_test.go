package batcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleInvocation(t *testing.T) {
	called := false
	b := NewBatcher(func() {
		called = true
	})

	b.Invoke()

	assert.Equal(t, true, called)
}

func TestTwoInvocations(t *testing.T) {
	count := 0
	move := make(chan int)
	b := NewBatcher(func() {
		<-move
		count++
	})

	go b.Invoke()
	go b.Invoke()
	move <- 0
	assert.Equal(t, 1, count)

}
