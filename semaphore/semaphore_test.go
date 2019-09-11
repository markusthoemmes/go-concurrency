// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.7

package semaphore

import (
	"context"
	"fmt"
	"testing"
)

func TestSemaphoreCorrectness(t *testing.T) {
	sem := newWaitListBasedSemaphore(1)
	sem.Acquire(context.Background())
	sem.Release()
}

// acquireN calls Acquire(size) on sem N times and then calls Release(size) N times.
func acquireN(b *testing.B, sem Interface, N int) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < N; j++ {
			sem.Acquire(context.Background())
		}
		for j := 0; j < N; j++ {
			sem.Release()
		}
	}
}

// tryAcquireN calls TryAcquire(size) on sem N times and then calls Release(size) N times.
func tryAcquireN(b *testing.B, sem Interface, N int) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < N; j++ {
			if !sem.TryAcquire() {
				b.Fatalf("TryAcquire() = false, want true")
			}
		}
		for j := 0; j < N; j++ {
			sem.Release()
		}
	}
}

func BenchmarkNewSeq(b *testing.B) {
	for _, cap := range []int64{1, 128} {
		b.Run(fmt.Sprintf("Weighted-%d", cap), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = newWaitListBasedSemaphore(cap)
			}
		})
		b.Run(fmt.Sprintf("semChan-%d", cap), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = newSemChan(cap)
			}
		})
	}
}

func BenchmarkAcquireSeq(b *testing.B) {
	for _, c := range []struct {
		cap, size int64
		N         int
	}{
		{1, 1, 1},
		{2, 1, 1},
		{16, 1, 1},
		{128, 1, 1},
		{2, 2, 1},
		{16, 2, 8},
		{128, 2, 64},
		{2, 1, 2},
		{16, 8, 2},
		{128, 64, 2},
	} {
		for _, w := range []struct {
			name string
			w    Interface
		}{
			{"new", newWaitListBasedSemaphore(c.cap)},
			{"semChan", newSemChan(c.cap)},
		} {
			b.Run(fmt.Sprintf("%s-acquire-%d-%d-%d", w.name, c.cap, c.size, c.N), func(b *testing.B) {
				acquireN(b, w.w, c.N)
			})
			b.Run(fmt.Sprintf("%s-tryAcquire-%d-%d-%d", w.name, c.cap, c.size, c.N), func(b *testing.B) {
				tryAcquireN(b, w.w, c.N)
			})
		}
	}
}
