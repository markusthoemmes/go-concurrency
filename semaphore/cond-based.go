package semaphore

import (
	"context"
	"sync"
	"sync/atomic"
)

func newCondBasedSemaphore(capacity int64) Interface {
	mux := sync.Mutex{}
	return &cond{capacity: capacity, cond: sync.NewCond(&mux)}
}

type cond struct {
	capacity int64
	cond     *sync.Cond

	// waitersCount keeps the count waiters.
	waitersCount int64
}

func (s *cond) Acquire(ctx context.Context) error {
	for !s.TryAcquire() {
		s.cond.L.Lock()
		// Try again under lock to make sure a racy release hasn't gotten in the way.
		if s.TryAcquire() {
			s.cond.L.Unlock()
			return ctx.Err()
		}
		// Increment the waitersCount **before** adding it to the list to make sure
		// Release acquires the lock and therefore sees the added waiter.
		atomic.AddInt64(&s.waitersCount, 1)
		s.cond.Wait()
		s.cond.L.Unlock()
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
	return ctx.Err()
}

func (s *cond) TryAcquire() bool {
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

func (s *cond) Release() {
	atomic.AddInt64(&s.capacity, 1)
	s.cond.Signal()
}
