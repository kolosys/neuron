# neuron Benchmarks

Performance benchmarks for the neuron package.

**Import Path:** `github.com/kolosys/neuron`

## Benchmark Results

The following benchmarks were run using `go test -bench=. -benchmem`:

| Benchmark | ns/op | B/op | allocs/op |
| --------- | ----- | ---- | --------- |
| `SerializeBody_Internal-4` | 170 | 56 | 2 |
| `Deduplicator_NoDedupe-4` | 530 | 400 | 7 |
| `Deduplicator_WithDedupe-4` | 342 | 384 | 6 |

### SerializeBody_Internal-4

- **Nanoseconds per operation:** 170 ns/op
- **Bytes allocated per operation:** 56 B/op
- **Allocations per operation:** 2 allocs/op
- **Number of runs:** 7082174

### Deduplicator_NoDedupe-4

- **Nanoseconds per operation:** 530 ns/op
- **Bytes allocated per operation:** 400 B/op
- **Allocations per operation:** 7 allocs/op
- **Number of runs:** 2268009

### Deduplicator_WithDedupe-4

- **Nanoseconds per operation:** 342 ns/op
- **Bytes allocated per operation:** 384 B/op
- **Allocations per operation:** 6 allocs/op
- **Number of runs:** 3469317

## Running Benchmarks

To run benchmarks for this package:

```bash
go test -bench=. -benchmem ./github.com/kolosys/neuron
```

To run benchmarks for all packages:

```bash
go test -bench=. -benchmem ./...
```

## External Links

- [Package Overview](../core-concepts/neuron.md)
- [API Reference](../api-reference/neuron.md)
- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/kolosys/neuron)
- [Source Code](https://github.com/kolosys/neuron/tree/dev/github.com/kolosys/neuron)
