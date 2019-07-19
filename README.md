# Golang Concurrency :lock: :rocket:

This repository is a playground repository to explore different semaphore (and other concurrency tools) implementations in Golang. 

The `sync` package leaves a lot to be desired in that regard so this aims as a place to gather different implementation mechanisms and unify them under a common interface. All implementations must pass the same test suite and can be run against the same benchmarks to get an understanding of which implementation is best suited for which use-cases or just simply superior to the others.

## Benchmarks

The current implementations of a gate yield the following results on my machine. It is worthwhile noting that only the channel based implementations actually fulfill the requirement to abort the `Wait` call as soon as the context is done.

### BenchmarkGateOpen

This benchmark simply checks the "happy" path. It measures the overhead the gate has if it is always open anyway.

```
goos: darwin
goarch: amd64
pkg: github.com/markusthoemmes/go-semaphores/gate
BenchmarkGateOpen/cond-based-1-12         	20000000	        62.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/cond-based-10-12        	20000000	        65.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/cond-based-100-12       	20000000	        72.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/cond-based-1000-12      	10000000	       503 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-atomic-based-1-12         	2000000000	         0.71 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-atomic-based-10-12        	2000000000	         0.73 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-atomic-based-100-12       	2000000000	         0.74 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-atomic-based-1000-12      	2000000000	         0.73 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-lock-based-1-12           	30000000	        51.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-lock-based-10-12          	30000000	        51.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-lock-based-100-12         	30000000	        51.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-lock-based-1000-12        	30000000	        52.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/lock-based-1-12                   	50000000	        40.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/lock-based-10-12                  	30000000	        43.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/lock-based-100-12                 	30000000	        44.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/lock-based-1000-12                	30000000	        46.6 ns/op	       0 B/op	       0 allocs/op
PASS
```

### BenchmarkGateOpenCloseFrequently

This benchmark measures what happens if we open and close the gate frequently during the test.

```
goos: darwin
goarch: amd64
pkg: github.com/markusthoemmes/go-semaphores/gate
BenchmarkGateOpenCloseFrequently/cond-based-1-12         	10000000	       153 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/cond-based-10-12        	10000000	       160 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/cond-based-100-12       	10000000	       150 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/cond-based-1000-12      	20000000	       190 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-atomic-based-1-12         	2000000000	         2.18 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-atomic-based-10-12        	2000000000	         0.83 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-atomic-based-100-12       	2000000000	         0.85 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-atomic-based-1000-12      	2000000000	         0.91 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-lock-based-1-12           	20000000	        80.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-lock-based-10-12          	20000000	        75.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-lock-based-100-12         	20000000	        65.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-lock-based-1000-12        	20000000	        54.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/lock-based-1-12                   	30000000	        55.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/lock-based-10-12                  	30000000	        47.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/lock-based-100-12                 	30000000	        48.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/lock-based-1000-12                	30000000	        47.6 ns/op	       0 B/op	       0 allocs/op
PASS
```