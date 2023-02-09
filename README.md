#### A thread-safe implement of an LRU cache in Go.
Optimised for concurrency


### Benchmarks
```
┌──[krypton] as k in ~/GolandProjects/logiq-cache on main
└──▶ go test -bench=.
goos: linux
goarch: amd64
pkg: logiq.ai/cache
cpu: Intel(R) Core(TM) i5-9300H CPU @ 2.40GHz
BenchmarkCache-8   	  456414	      3527 ns/op
PASS
ok  	logiq.ai/cache	1.707s
```

### Tests
Test cases and benchmarks are situated at main_test.go