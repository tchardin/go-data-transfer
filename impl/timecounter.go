package impl

import (
	"sync/atomic"
	"time"
)

// timeCounter is used to generate a monotonically increasing sequence.
// It starts at the current time, then increments on each call to next.
type timeCounter struct {
	counter uint64
}

func newTimeCounter() *timeCounter {
	return &timeCounter{counter: uint64(time.Now().UnixNano())}
}

func (tc *timeCounter) next() uint64 {
	counter := atomic.AddUint64(&tc.counter, 1)
	return counter
}
