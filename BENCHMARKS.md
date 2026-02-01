# Benchmarks

This file tracks baseline benchmark results so changes can be compared over time.
Results vary by machine, OS, and Go version.

## 2026-02-01 baseline

Command:

```bash
go test -bench . -benchmem -run ^$
```

Environment:

- go version go1.25.2 darwin/arm64
- CPU: Apple M1

Results:

```
BenchmarkGreatCircleDistance-8             31560385        37.79 ns/op        0 B/op        0 allocs/op
BenchmarkRhumbLineDistance-8               46364043        25.93 ns/op        0 B/op        0 allocs/op
BenchmarkGreatCircleIntermediatePoint-8     7621258       155.9 ns/op        0 B/op        0 allocs/op
BenchmarkGreatCirclePointAtSpeed-8          6016376       197.5 ns/op        0 B/op        0 allocs/op
BenchmarkGreatCircleProject-8               3919816       305.2 ns/op        0 B/op        0 allocs/op
BenchmarkGreatCircleProjectToSegment-8      3559207       335.3 ns/op        0 B/op        0 allocs/op
BenchmarkGeohashEncode-8                    9563523       124.0 ns/op       24 B/op        2 allocs/op
BenchmarkGeohashDecode-8                   12052099        99.54 ns/op       0 B/op        0 allocs/op
BenchmarkGeohashNeighbors-8                  977490      1229 ns/op        192 B/op       16 allocs/op
BenchmarkDijkstra-8                           39055     30204 ns/op       41505 B/op     1007 allocs/op
BenchmarkTSPNearestNeighbor-8               5058968       236.6 ns/op       280 B/op        6 allocs/op
```

## 2026-02-01 after geohash preallocation

Command:

```bash
go test -bench . -benchmem -run ^$
```

Environment:

- go version go1.25.2 darwin/arm64
- CPU: Apple M1

Results:

```
BenchmarkGreatCircleDistance-8             31574508        37.78 ns/op        0 B/op        0 allocs/op
BenchmarkRhumbLineDistance-8               46692302        26.02 ns/op        0 B/op        0 allocs/op
BenchmarkGreatCircleIntermediatePoint-8     7593616       157.6 ns/op        0 B/op        0 allocs/op
BenchmarkGreatCirclePointAtSpeed-8          6010648       199.1 ns/op        0 B/op        0 allocs/op
BenchmarkGreatCircleProject-8               3927168       304.6 ns/op        0 B/op        0 allocs/op
BenchmarkGreatCircleProjectToSegment-8      3580462       334.9 ns/op        0 B/op        0 allocs/op
BenchmarkGeohashEncode-8                   11390935       104.1 ns/op       16 B/op        1 allocs/op
BenchmarkGeohashDecode-8                   11981392        99.56 ns/op       0 B/op        0 allocs/op
BenchmarkGeohashNeighbors-8                 1000000      1222 ns/op        128 B/op        8 allocs/op
BenchmarkDijkstra-8                           39348     30321 ns/op       41505 B/op     1007 allocs/op
BenchmarkTSPNearestNeighbor-8               5057041       237.2 ns/op       280 B/op        6 allocs/op
```
