package main

import "sync"

type ErrorCounter struct {
	sync.RWMutex
	value int
}

func (counter *ErrorCounter) Increase() {
	counter.Lock()
	counter.value++
	counter.Unlock()
}

func (counter *ErrorCounter) Less(i int) bool {
	counter.RLock()
	defer counter.RUnlock()
	return counter.value < i
}

func (counter *ErrorCounter) Value() int {
	counter.RLock()
	defer counter.RUnlock()
	return counter.value
}
