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
	"sync/atomic"
	"unsafe"
)

// newChannelBasedWithPointers returns a gate that is based on a broadcast
// channel wrapped in an unsafe.Pointer, which acts as state.
// Once the gate closes a channel is generated which concurrent Wait calls
// read from. Once the gate opens again, that channel is closed, notifying
// all waiters to move on.
func newChannelBasedWithPointers() Interface {
	initial := make(chan struct{})
	return &channelBasedWithPointers{
		broadcast: unsafe.Pointer(&initial),
	}
}

type channelBasedWithPointers struct {
	broadcast unsafe.Pointer
}

// Open implements interface.
func (c *channelBasedWithPointers) Open() {
	channelPtr := atomic.SwapPointer(&c.broadcast, nil)
	if channelPtr != nil {
		channel := *(*chan struct{})(channelPtr)
		close(channel)
	}
}

// Close implements interface.
func (c *channelBasedWithPointers) Close() {
	channel := make(chan struct{})
	atomic.CompareAndSwapPointer(&c.broadcast, nil, unsafe.Pointer(&channel))
}

// Wait implements interface.
func (c *channelBasedWithPointers) Wait(ctx context.Context) error {
	channelPtr := atomic.LoadPointer(&c.broadcast)
	if channelPtr == nil {
		return ctx.Err()
	}

	channel := *(*chan struct{})(channelPtr)
	select {
	case <-channel:
		return ctx.Err()
	case <-ctx.Done():
		return ctx.Err()
	}
}
