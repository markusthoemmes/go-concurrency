package gate

import (
	"context"
	"sync"
)

// newLockBased returns a gate that is based on sync.RWMutex.
// The gate is open when the lock is unlocked, close if it locked.
// Concurrent waits are possible via only acquiring read locks.
func newLockBased() Interface {
	gate := &lockBased{}
	gate.mux.Lock()
	return gate
}

type lockBased struct {
	mux sync.RWMutex
}

func (c *lockBased) Open() {
	c.mux.Unlock()
}

func (c *lockBased) Close() {
	c.mux.Lock()
}

func (c *lockBased) Wait(ctx context.Context) error {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return ctx.Err()
}
