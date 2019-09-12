// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.7

package semaphore

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var (
	implementations = map[string]func(int64) Interface{
		"channel":   newChannelBased,
		"wait-list": newWaitListBasedSemaphore,
		"cond":      newCondBasedSemaphore,
	}

	parallelisms = []int{1, 10, 100, 1000}

	bg = context.Background()
)

func TestSemaphoreCorrectness(t *testing.T) {
	for implName, impl := range implementations {
		t.Run(implName, func(t *testing.T) {
			sem := impl(1)
			sem.Acquire(bg)
			sem.Release()
		})
	}
}

func BenchmarkSemaphore(b *testing.B) {
	for implName, impl := range implementations {
		for _, capacity := range []int64{1, 8, 128} {
			for _, parallelism := range parallelisms {
				benchName := fmt.Sprintf("%s-capacity:%d-threads:%d", implName, capacity, parallelism)
				b.Run(benchName, func(b *testing.B) {
					b.SetParallelism(parallelism)
					s := impl(capacity)

					b.RunParallel(func(pb *testing.PB) {
						for pb.Next() {
							s.Acquire(bg)
							// Actually create some contention.
							time.Sleep(10 * time.Nanosecond)
							s.Release()
						}
					})
				})
			}
		}
	}
}
