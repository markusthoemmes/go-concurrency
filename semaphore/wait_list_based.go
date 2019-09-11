package semaphore

import (
	"container/list"
	"context"
	"sync"
	"sync/atomic"
)

func newWaitListBasedSemaphore(capacity int64) Interface {
	return &sema{capacity: capacity}
}

type sema struct {
	capacity int64

	waitersMux   sync.RWMutex
	waiters      list.List
	waitersCount int64
}

func (s *sema) Acquire(ctx context.Context) error {
	firstRun := true
	for !s.TryAcquire() {
		s.waitersMux.Lock()
		// Try again under lock to make sure a racy release hasn't gotten in the way.
		if s.TryAcquire() {
			s.waitersMux.Unlock()
			return ctx.Err()
		}

		ready := make(chan struct{})
		var elem *list.Element
		// Increment the waitersCount **before** adding it to the list to make sure
		// Release acquires the lock and therefore sees the added waiter.
		atomic.AddInt64(&s.waitersCount, 1)
		if firstRun {
			// Push the channel to the back of the waiter queue on the first pass.
			elem = s.waiters.PushBack(ready)
			firstRun = false
		} else {
			// Push the channel to the front of the queue to prevent starvation.
			elem = s.waiters.PushFront(ready)
		}
		s.waitersMux.Unlock()

		select {
		case <-ctx.Done():
			s.waitersMux.Lock()
			s.waiters.Remove(elem)
			s.waitersMux.Unlock()
			// Bail out early.
			return ctx.Err()
		case <-ready:
		}
	}
	return ctx.Err()
}

func (s *sema) TryAcquire() bool {
	for {
		cur := atomic.LoadInt64(&s.capacity)
		if cur == 0 {
			return false
		}
		if atomic.CompareAndSwapInt64(&s.capacity, cur, cur-1) {
			return true
		}
	}
}

func (s *sema) Release() {
	atomic.AddInt64(&s.capacity, 1)
	if atomic.LoadInt64(&s.waitersCount) == 0 {
		return
	}

	s.waitersMux.Lock()
	defer s.waitersMux.Unlock()

	elem := s.waiters.Front()
	if elem == nil {
		return
	}

	s.waiters.Remove(elem)
	// Decrement the waitersCount **after** removing it. Worst case is we see
	// a nil and bail out early.
	atomic.AddInt64(&s.waitersCount, -1)
	close(elem.Value.(chan struct{}))
}
