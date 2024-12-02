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

// RunWithResult executes a function asynchronously and returns the result via a channel
func (a *Async) RunWithResult(fn func() (interface{}, error)) <-chan Result {
	a.wg.Add(1)
	resultChan := make(chan Result, 1)
	go func() {
		defer a.wg.Done()
		defer close(resultChan)

		result, err := fn()
		resultChan <- Result{Value: result, Err: err}
	}()
	return resultChan
}

func (a *Async) Wait() {
	a.wg.Wait()
}

type Result struct {
	Value interface{}
	Err   error
}
