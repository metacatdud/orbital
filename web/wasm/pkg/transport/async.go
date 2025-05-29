package transport

import (
	"sync"
)

type Async struct {
	wg sync.WaitGroup
}

// Async executes a function asynchronously without returning a value simmilar to a JS callback
func (a *Async) Async(fn func()) {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		fn()
	}()
}

// Wait for all async executions to finish.
// Use this only if parallel execution of more than two async calls is required
func (a *Async) Wait() {
	a.wg.Wait()
}
