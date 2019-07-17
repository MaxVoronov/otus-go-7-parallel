package main

import "sync"

type ErrorCounter struct {
	sync.RWMutex
	Value int
}

func (counter *ErrorCounter) Increase() {
	counter.Lock()
	counter.Value++
	counter.Unlock()
}

func (counter *ErrorCounter) Less(i int) bool {
	counter.RLock()
	defer counter.RUnlock()
	return counter.Value < i
}
