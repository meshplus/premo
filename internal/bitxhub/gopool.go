package bitxhub

import "sync"

type Pool struct {
	queue chan struct{}
	wg    *sync.WaitGroup
}

func NewGoPool(size int) *Pool {
	if size <= 0 {
		size = 1
	}
	return &Pool{
		queue: make(chan struct{}, size),
		wg:    &sync.WaitGroup{},
	}
}

func (p *Pool) Add() {
	p.queue <- struct{}{}
	p.wg.Add(1)
}

func (p *Pool) Done() {
	<-p.queue
	p.wg.Done()
}

func (p *Pool) Wait() {
	p.wg.Wait()
}
