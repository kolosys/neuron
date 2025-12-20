# neuron Benchmarks

Performance benchmarks for the neuron package.

**Import Path:** `github.com/kolosys/neuron`

## Benchmark Results

The following benchmarks were run using `go test -bench=. -benchmem`:

| Benchmark | ns/op | B/op | allocs/op |
| --------- | ----- | ---- | --------- |
| `SerializeBody_Internal-4` | 210 | 56 | 2 |

### SerializeBody_Internal-4

- **Nanoseconds per operation:** 210 ns/op
- **Bytes allocated per operation:** 56 B/op
- **Allocations per operation:** 2 allocs/op
- **Number of runs:** 5294348

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
