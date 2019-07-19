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
BenchmarkGateOpen/channel-lock-1-12         	30000000	        41.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-lock-10-12        	30000000	        41.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-lock-100-12       	30000000	        41.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-lock-1000-12      	30000000	        42.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/lock-1-12                 	50000000	        38.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/lock-10-12                	50000000	        38.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/lock-100-12               	50000000	        38.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/lock-1000-12              	50000000	        38.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/cond-1-12                 	20000000	        64.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/cond-10-12                	20000000	        64.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/cond-100-12               	20000000	        72.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/cond-1000-12              	20000000	       560 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-atomics-1-12      	2000000000	         0.77 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-atomics-10-12     	2000000000	         0.77 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-atomics-100-12    	2000000000	         0.76 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-atomics-1000-12   	2000000000	         0.78 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-pointers-1-12     	2000000000	         0.78 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-pointers-10-12    	2000000000	         0.81 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-pointers-100-12   	2000000000	         0.81 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-pointers-1000-12  	2000000000	         0.81 ns/op	       0 B/op	       0 allocs/op
PASS
```

### BenchmarkGateOpenCloseFrequently

This benchmark measures what happens if we open and close the gate frequently during the test.

```
goos: darwin
goarch: amd64
pkg: github.com/markusthoemmes/go-semaphores/gate
BenchmarkGateOpenCloseFrequently/channel-atomics-1-12         	2000000000	         1.43 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-atomics-10-12        	2000000000	         0.80 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-atomics-100-12       	2000000000	         0.73 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-atomics-1000-12      	2000000000	         0.75 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-pointers-1-12        	500000000	         2.53 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-pointers-10-12       	2000000000	         0.92 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-pointers-100-12      	2000000000	         0.88 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-pointers-1000-12     	2000000000	         0.83 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-lock-1-12            	 5000000	       291 ns/op	     165 B/op	       2 allocs/op
BenchmarkGateOpenCloseFrequently/channel-lock-10-12           	20000000	       136 ns/op	      56 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-lock-100-12          	30000000	        43.2 ns/op	       1 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-lock-1000-12         	20000000	        56.3 ns/op	       7 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/lock-1-12                    	10000000	       154 ns/op	      67 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/lock-10-12                   	10000000	       115 ns/op	      37 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/lock-100-12                  	20000000	        51.5 ns/op	       2 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/lock-1000-12                 	30000000	        54.8 ns/op	       2 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/cond-1-12                    	20000000	       112 ns/op	      35 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/cond-10-12                   	20000000	        86.4 ns/op	      23 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/cond-100-12                  	20000000	       100.0 ns/op	      24 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/cond-1000-12                 	10000000	       161 ns/op	      32 B/op	       0 allocs/op
PASS
```