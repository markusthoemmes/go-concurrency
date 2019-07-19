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
BenchmarkGateOpen/cond-based-1-12         	20000000	        63.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/cond-based-10-12        	20000000	        65.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/cond-based-100-12       	20000000	        72.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/cond-based-1000-12      	10000000	       527 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-atomic-based-1-12         	2000000000	         0.68 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-atomic-based-10-12        	2000000000	         0.77 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-atomic-based-100-12       	2000000000	         0.72 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-atomic-based-1000-12      	2000000000	         0.79 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-lock-based-1-12           	20000000	        55.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-lock-based-10-12          	30000000	        51.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-lock-based-100-12         	30000000	        51.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-lock-based-1000-12        	30000000	        51.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/lock-based-1-12                   	50000000	        36.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/lock-based-10-12                  	50000000	        35.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/lock-based-100-12                 	50000000	        35.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/lock-based-1000-12                	50000000	        35.5 ns/op	       0 B/op	       0 allocs/op
PASS
```

### BenchmarkGateOpenCloseFrequently

This benchmark measures what happens if we open and close the gate frequently during the test.

```
goos: darwin
goarch: amd64
pkg: github.com/markusthoemmes/go-semaphores/gate
BenchmarkGateOpenCloseFrequently/cond-based-1-12         	20000000	        96.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/cond-based-10-12        	20000000	       155 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/cond-based-100-12       	20000000	       113 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/cond-based-1000-12      	10000000	       142 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-atomic-based-1-12         	1000000000	         2.60 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-atomic-based-10-12        	1000000000	         1.00 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-atomic-based-100-12       	1000000000	         0.90 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-atomic-based-1000-12      	2000000000	         0.85 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-lock-based-1-12           	 5000000	       240 ns/op	      53 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-lock-based-10-12          	20000000	        85.3 ns/op	      10 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-lock-based-100-12         	20000000	        59.3 ns/op	       1 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-lock-based-1000-12        	30000000	        42.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/lock-based-1-12                   	20000000	       121 ns/op	      26 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/lock-based-10-12                  	30000000	        92.5 ns/op	      17 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/lock-based-100-12                 	20000000	        70.1 ns/op	       1 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/lock-based-1000-12                	20000000	        53.7 ns/op	       2 B/op	       0 allocs/op
PASS
```