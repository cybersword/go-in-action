// Package work manages a pool of goroutines to perform work.
package work

import (
	"sync"
)

// Worker must be implemented by tpyes that want to use
// the work pool.
type Worker interface {
	Task()
}

// Pool provides a pool of goroutines that can execute
// any Worker tasks that are submitted.
type Pool struct {
	work chan Worker
	wg   sync.WaitGroup
}

// New creates a new work pool.
func New(maxGoroutines int) *Pool {
	p := Pool{
		work: make(chan Worker),
	}

	p.wg.Add(maxGoroutines)
	// 后台起多个worker(goroutine), 每个都从chan等待数据
	for i := 0; i < maxGoroutines; i++ {
		go func() {
			// 只有chan关闭,并为空,才会退出循环
			for w := range p.work {
				w.Task()
			}
			p.wg.Done()
		}()
	}

	return &p
}

// Run submits work to the pool.
func (p *Pool) Run(w Worker) {
	p.work <- w
}

// Shutdown waits for all the goroutines to shutdown.
func (p *Pool) Shutdown() {
	close(p.work)
	p.wg.Wait()
}
