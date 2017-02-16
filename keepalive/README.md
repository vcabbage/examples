# Example of argument liveness change in Go 1.8

> The garbage collector no longer considers arguments live throughout the entirety of a function.
> For more information, and for how to force a variable to remain live, see the runtime.KeepAlive
> function added in Go 1.7.
https://beta.golang.org/doc/go1.8#liveness

### Description

This repo points out an issue which can occur when passing Go allocated data to C and never
directly referencing the Go data again.

1. The example `DoIt()` function takes a Go slice of `int64` and passes the underlying array to
C.
1. The C code passes back a pointer to a C struct.
1. To simulate work being done, GC is forced and a slice of random size is created and filled
with random numbers.
1. The C struct is passed to a C function which sums the values in the contained array.
1. `DoIt()` returns the result.

`DoItKeepAlive()` performs the exact same operation except it inserts a call to `runtime.KeepAlive()`
to ensure the garbage collector does not collect the slice of `int64`.

## Output

Go 1.7:
```
# go1.7 version
go version go1.7.5 darwin/amd64
# go1.7 test -v
=== RUN   TestDoIt
--- PASS: TestDoIt (0.00s)
=== RUN   TestDoItKeepAlive
--- PASS: TestDoItKeepAlive (0.00s)
PASS
ok  	github.com/vcabbage/examples/keepalive	0.009s
# go1.7 test -v
=== RUN   TestDoIt
--- PASS: TestDoIt (0.00s)
=== RUN   TestDoItKeepAlive
--- PASS: TestDoItKeepAlive (0.00s)
PASS
ok  	github.com/vcabbage/examples/keepalive	0.009s
# go1.7 test -v
=== RUN   TestDoIt
--- PASS: TestDoIt (0.00s)
=== RUN   TestDoItKeepAlive
--- PASS: TestDoItKeepAlive (0.00s)
PASS
ok  	github.com/vcabbage/examples/keepalive	0.009s
```

Go 1.8:
```
# go1.8 version
go version go1.8 darwin/amd64
# go1.8 test -v
=== RUN   TestDoIt
--- FAIL: TestDoIt (0.00s)
	keepalive_test.go:18: Wanted 312487500, got 2459676341428310292
=== RUN   TestDoItKeepAlive
--- PASS: TestDoItKeepAlive (0.00s)
FAIL
exit status 1
FAIL	github.com/vcabbage/examples/keepalive	0.009s
# go1.8 test -v
=== RUN   TestDoIt
--- FAIL: TestDoIt (0.00s)
	keepalive_test.go:18: Wanted 312487500, got 2888024562982026679
=== RUN   TestDoItKeepAlive
--- PASS: TestDoItKeepAlive (0.00s)
FAIL
exit status 1
FAIL	github.com/vcabbage/examples/keepalive	0.009s
# go1.8 test -v
=== RUN   TestDoIt
--- PASS: TestDoIt (0.00s)
=== RUN   TestDoItKeepAlive
--- PASS: TestDoItKeepAlive (0.00s)
PASS
ok  	github.com/vcabbage/examples/keepalive	0.008s
```
As shown, the behavior is not consistent between runs. The inconsistency is made more apparent by
manually running GC and creating a slice of a random size, simulating an additional amount of allocations
that are different beteween runs.