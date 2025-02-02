[![Go Report Card](https://goreportcard.com/badge/github.com/liyiheng/zset)](https://goreportcard.com/report/github.com/liyiheng/zset)
# zset
Implementing sorted set in Redis with golang.

## Installation
```bash
go get -u github.com/liyiheng/zset
```

## Usage

```go
s := zset.New[int64]()
// add data
s.Set(66, 1001)
s.Set(77, 1002)
s.Set(88, 1003)
s.Set(100, 1004)
s.Set(99, 1005)
s.Set(44, 1006)
// update data
s.Set(44, 1001)

// get rank by key
rank, score := s.GetRank(1004, false)
// get data by rank
id, score := s.GetDataByRank(0, true)

s.Delete(1001)

// Increase score
s.IncrBy(5.0, 1001)

// ZRANGE, ASC
five := make([]int64, 0, 5)
s.Range(0, 5, func(score float64, k int64) {
	five = append(five, k)
})

// ZREVRANGE, DESC
all := make([]int64, 0)
s.RevRange(0, -1, func(score float64, k int64) {
	all = append(all, k)
})


```

## Benchmark

```text
 OS: Arch Linux 
 Kernel: x86_64 Linux 5.1.5-arch1-2-ARCH
 CPU: Intel Core i7-8750H @ 12x 4.1GHz [46.0°C]
 RAM: 3295MiB / 7821MiB
```

```bash
go test -test.bench=".*"
goos: linux
goarch: amd64
BenchmarkSortedSet_Add-12              	 1000000	      3050 ns/op
BenchmarkSortedSet_GetRank-12          	  500000	      2963 ns/op
BenchmarkSortedSet_GetDataByRank-12    	 2000000	       620 ns/op
PASS
```
