#### A thread-safe implement of an LRU cache in Go.
Optimised for concurrency


### Benchmarks
```
go test -bench=.
goos: linux
goarch: amd64
pkg: logiq.ai/cache
BenchmarkCache-8   	  503358	      4329 ns/op

PASS
ok  	logiq.ai/cache	2.274s
```

### Tests
Test cases and benchmarks are situated at main_test.go