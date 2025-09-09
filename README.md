# Consistent Hashing Tests

## Purpose

This laboratory work demonstrates the implementation of a consistent hashing algorithm in Go. The project includes:

- A custom implementation of consistent hashing using CRC32 hashing.
- Comparison tests with two popular Go libraries: `github.com/ArchishmanSengupta/consistent-hashing` and `github.com/buraksezer/consistent`.
- Benchmarking to measure performance.

The goal is to understand how consistent hashing works, its benefits for distributed systems, and how different implementations compare in terms of load distribution and performance.

## Implementation

The main implementation is in [`main.go`](main.go), which defines the `ConsistentHash` struct with methods:

- `NewConsistentHash`: Creates a new instance with configurable replicas and hash function.
- `Add`: Adds nodes to the hash ring with virtual replicas.
- `Get`: Retrieves the node responsible for a given key.

## Testing

Tests in [`main_test.go`](main_test.go) include:

- Basic functionality test for the custom implementation.
- Comparison with ArchishmanSengupta's library.
- Comparison with Buraksezer's library, including node removal scenarios.
- Load distribution analysis.

## Benchmarking

Benchmarks compare the performance of the custom implementation and ArchishmanSengupta's library.

## Dependencies

- Go 1.24.3
- github.com/ArchishmanSengupta/consistent-hashing v1.0.2
- github.com/cespare/xxhash v1.1.0
- github.com/buraksezer/consistent v0.10.0

## Running the Tests

Run `go test` to execute the tests and benchmarks.

Run `go test -bench=.` to run benchmarks.