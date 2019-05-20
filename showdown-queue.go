package main

import (
	"sync"
)

type ShowdownQueue struct {
	mut          sync.Mutex
	waitGroup    sync.WaitGroup
	symbols      []string
	currentIndex int
}

func (q *ShowdownQueue) SetSymbols(array []string) {
	q.symbols = array
}

func (q *ShowdownQueue) GetNextSymbol() string {
	var res string
	q.mut.Lock()
	if q.currentIndex < len(q.symbols) {
		res = q.symbols[q.currentIndex]
		q.currentIndex += 1
	} else {
		res = ""
	}
	q.mut.Unlock()
	return res
}
