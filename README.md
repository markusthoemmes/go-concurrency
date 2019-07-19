# Golang Concurrency :lock: :rocket:

This repository is a playground repository to explore different semaphore (and other concurrency tools) implementations in Golang. 

The `sync` package leaves a lot to be desired in that regard so this aims as a place to gather different implementation mechanisms and unify them under a common interface. All implementations must pass the same test suite and can be run against the same benchmarks to get an understanding of which implementation is best suited for which use-cases or just simply superior to the others.

## Benchmarks

The current implementations of a gate yield the following results on my machine

### BenchmarkGateOpen

This benchmark simply checks the "happy" path. It measures the overhead the gate has if it is always open anyway.

```
goos: darwin
goarch: amd64
pkg: github.com/markusthoemmes/go-semaphores/gate
BenchmarkGateOpen/cond-based-12         	20000000	        62.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpen/channel-based-12      	2000000000	         0.67 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/markusthoemmes/go-semaphores/gate	3.242s
Success: Benchmarks passed.
```

### BenchmarkGateOpenCloseFrequently

This benchmark measures what happens if we open and close the gate frequently during the test.

```
goos: darwin
goarch: amd64
pkg: github.com/markusthoemmes/go-semaphores/gate
BenchmarkGateOpenCloseFrequently/cond-based-12         	10000000	       153 ns/op	       0 B/op	       0 allocs/op
BenchmarkGateOpenCloseFrequently/channel-based-12      	2000000000	         1.86 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/markusthoemmes/go-semaphores/gate	5.636s
Success: Benchmarks passed.
```