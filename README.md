# Regex Parser

## Motivation

Simple regex parser written to understand

- Lexers.
- Parsers.
- NFA / DFA.
- How to represent them.
- Graph traversal.
- Go Error handling / Folder structure.
- Node parsing.

## Benchmarking

My implementation vs native implementation for something like checking if emails match the regex is over `31x` slower.

```
go test -bench=. -benchtime=100000x ./tests/perf/
```

```
  󱞪 time go test -bench=. -benchtime=100000x ./tests/perf/test_naive_test.go
goos: linux
goarch: amd64
cpu: AMD Ryzen 7 5800H with Radeon Graphics
BenchmarkN-16    	  100000	     77874 ns/op
PASS
ok  	command-line-arguments	7.791s
go test -bench=. -benchtime=100000x ./tests/perf/test_naive_test.go  8.10s user 0.32s system 106% cpu 7.914 total

---

  󱞪 time go test -bench=. -benchtime=100000x ./tests/perf/test_inbuilt_test.go
goos: linux
goarch: amd64
cpu: AMD Ryzen 7 5800H with Radeon Graphics
BenchmarkGo-16    	  100000	      2603 ns/op
PASS
ok  	command-line-arguments	0.264s
go test -bench=. -benchtime=100000x ./tests/perf/test_inbuilt_test.go  0.43s user 0.11s system 139% cpu 0.387 total
```

Now i will run some tests and figure out how to optimize this.

## Changelogs

I implemented bitmask to speed up `[]` syntax. Because i used to generate all the values between the brackets and then connect them to multiple NFA objects. Instead i can have a single connection which is a bitmap which will tell me what values can pass through. Using simple math.

`[4]uint64 bitMap`

Then i can do something like

```
for ch := range characters {
    bitMap[ch/64] |= 1 << (ch%64)
}
```

This will set bits for characters that are possible.

Then i can check in `compile.go` sometihng like

```
if bitMap[ch/64] & (1 << (ch % 64)) {}
```

Then its found and i can add the state to NFA.

```
  󱞪 time go test -bench=. -benchtime=100000x ./tests/perf/test_naive_test.go
goos: linux
goarch: amd64
cpu: AMD Ryzen 7 5800H with Radeon Graphics
BenchmarkN-16    	  100000	     55975 ns/op
PASS
ok  	command-line-arguments	5.601s
go test -bench=. -benchtime=100000x ./tests/perf/test_naive_test.go  5.87s user 0.36s system 108% cpu 5.729 total
```

We shaved off around 2s!!!!!

---
**No LLMS were used**
