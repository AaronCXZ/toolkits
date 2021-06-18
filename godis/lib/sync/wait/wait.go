package wait

import (
	"sync"
	"time"
)

type Wait struct {
	wg sync.WaitGroup
}

// Add WaitGroup的Add方法
func (w *Wait) Add(delta int) {
	w.wg.Add(delta)
}

// Done WaitGroup的Done方法
func (w *Wait) Done() {
	w.wg.Done()
}

// Wait WaitGroup的Wait方法
func (w *Wait) Wait() {
	w.wg.Wait()
}

// WaitWithTimeout timeout是返回true
func (w *Wait) WaitWithTimeout(timeout time.Duration) bool {
	c := make(chan bool)
	go func() {
		defer close(c)
		w.wg.Wait()
		c <- true
	}()
	select {
	case <-c:
		return false
	case <-time.After(timeout):
		return true
	}
}
