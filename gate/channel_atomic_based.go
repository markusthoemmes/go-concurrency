package gate

import (
	"context"
	"sync"
	"sync/atomic"
)

// newChannelBasedWithAtomics returns a gate that is based on a broadcast
// channel and an atomically updated open/closed state, used as a fast-path.
// Once the gate closes a channel is generated which concurrent Wait calls
// read from. Once the gate opens again, that channel is closed, notifying
// all waiters to move on.
func newChannelBasedWithAtomics() Interface {
	return &channelBasedWithAtomics{
		broadcast: make(chan struct{}),
	}
}

type channelBasedWithAtomics struct {
	broadcast chan struct{}
	open      uint32

	// mux guards Open and Close from being called concurrently.
	mux sync.Mutex
}

func (c *channelBasedWithAtomics) Open() {
	c.mux.Lock()
	defer c.mux.Unlock()
	if atomic.CompareAndSwapUint32(&c.open, 0, 1) {
		close(c.broadcast)
	}
}

func (c *channelBasedWithAtomics) Close() {
	c.mux.Lock()
	defer c.mux.Unlock()
	if atomic.CompareAndSwapUint32(&c.open, 1, 0) {
		c.broadcast = make(chan struct{})
	}
}

func (c *channelBasedWithAtomics) Wait(ctx context.Context) error {
	if atomic.LoadUint32(&c.open) > 0 {
		return ctx.Err()
	}
	select {
	case <-c.broadcast:
		return ctx.Err()
	case <-ctx.Done():
		return ctx.Err()
	}
}
