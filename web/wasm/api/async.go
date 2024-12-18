// Package api Async exists to wrap the underlying fetch function which cause the not handled
// promise block the main event loop block and an immediate deadlock
package api

import (
	"sync"
)

type Async struct {
	wg sync.WaitGroup
}

func NewAsync() *Async {
	return &Async{}
}

// Run executes a function asynchronously without returning a value
func (a *Async) Run(fn func()) {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		fn()
	}()
}

func (a *Async) Wait() {
	a.wg.Wait()
}
