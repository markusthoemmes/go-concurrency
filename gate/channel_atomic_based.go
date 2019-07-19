/*
Copyright 2019 Markus Thömmes

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

	// mux guards Open and Close from being called concurrently. It also
	// guards broadcast state changes.
	mux sync.RWMutex
}

// Open implements interface.
func (c *channelBasedWithAtomics) Open() {
	c.mux.Lock()
	defer c.mux.Unlock()
	if atomic.CompareAndSwapUint32(&c.open, 0, 1) {
		close(c.broadcast)
	}
}

// Close implements interface.
func (c *channelBasedWithAtomics) Close() {
	c.mux.Lock()
	defer c.mux.Unlock()
	if atomic.CompareAndSwapUint32(&c.open, 1, 0) {
		c.broadcast = make(chan struct{})
	}
}

// Wait implements interface.
func (c *channelBasedWithAtomics) Wait(ctx context.Context) error {
	if atomic.LoadUint32(&c.open) > 0 {
		return ctx.Err()
	}
	select {
	case <-c.broadcaster():
		return ctx.Err()
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *channelBasedWithAtomics) broadcaster() chan struct{} {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.broadcast
}
