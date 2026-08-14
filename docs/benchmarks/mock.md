# mock Benchmarks

Performance benchmarks for the mock package.

**Import Path:** `github.com/kolosys/neuron/mock`

## No Benchmarks Available

No benchmark results are available for this package. To add benchmarks:

1. Create a `*_test.go` file in the package directory
2. Add benchmark functions following the pattern:
   ```go
   func BenchmarkFunctionName(b *testing.B) {
       for i := 0; i < b.N; i++ {
           // Your code here
       }
   }
   ```
3. Run `proton benchmark` to generate benchmark results

## Running Benchmarks

To run benchmarks for this package:

```bash
go test -bench=. -benchmem ./mock
```

To run benchmarks for all packages:

```bash
go test -bench=. -benchmem ./...
```

## External Links

- [Package Overview](../core-concepts/mock.md)
- [API Reference](../api-reference/mock.md)
- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/kolosys/neuron/mock)
- [Source Code](https://github.com/kolosys/neuron/tree/dev/mock)
