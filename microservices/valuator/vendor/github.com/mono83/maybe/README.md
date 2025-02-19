Maybe
=====

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/mono83/maybe.svg)](https://github.com/mono83/maybe)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/mono83/maybe)
[![Go Report Card](https://goreportcard.com/badge/github.com/mono83/maybe)](https://goreportcard.com/report/github.com/mono83/maybe)
[![CircleCI](https://dl.circleci.com/status-badge/img/gh/mono83/maybe/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/mono83/maybe/tree/main)

Simple zero dependency struct-based generic monad implementation in Go.

## Installation

```bash
go get -u github.com/mono83/maybe
```

## Usage

### Instantiation

There are two main ways to create `maybe.Maybe` monadic container - it is either call to `maybe.Nothing()` or 
`maybe.Just(value)`. First will create empty monad, while second one will produce not empty monad with given value. 

Pay attention, that `maybe.Just` does not verify content of passed value, so constructions like 
`maybe.Just[error](nil)` can create not empty monads with nil value within. To avoid this behaviour use 
`maybe.Nilable`, that will perform `nil` check of given value.

In addition to raw instantiation, this library provide mapping constructor `maybe.Map`, able to convert 
`Maybe` from one type to another using mapping function

```go
n := maybe.Nothing[string]()    // String type, but empty
i := maybe.Just(1)              // Int type
s := maybe.Map(i, strconv.Itoa) // String type
```

### JSON

JSON bindings can be accessed using `github.com/mono83/maybe/json` package:

```go
import "github.com/mono83/maybe/json"

type Request struct {
	ID       int
	ParentID json.Maybe[int]
}

```

### Benchmarks

```
$ go test -bench=. -benchmem
goos: linux
goarch: amd64
pkg: github.com/mono83/maybe
cpu: Intel(R) Core(TM) i5-10400F CPU @ 2.90GHz
BenchmarkNothing-12             1000000000               0.2704 ns/op          0 B/op          0 allocs/op
BenchmarkJust-12                904382658                1.313 ns/op           0 B/op          0 allocs/op
BenchmarkPtr-12                 1000000000               1.051 ns/op           0 B/op          0 allocs/op
BenchmarkPtrNil-12              755984685                1.578 ns/op           0 B/op          0 allocs/op
BenchmarkNilableInt-12          72761029                16.37 ns/op            8 B/op          0 allocs/op
BenchmarkNilableNil-12          275986280                4.157 ns/op           0 B/op          0 allocs/op
PASS
ok      github.com/mono83/maybe 6.931s
```