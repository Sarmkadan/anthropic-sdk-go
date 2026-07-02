# Performance Benchmarks for anthropic-sdk-go

This directory contains performance benchmarks for the anthropic-sdk-go library using [BenchmarkDotNet](https://benchmarkdotnet.org/).


## Running Benchmarks

### Prerequisites

- Go 1.23.0 or later
- BenchmarkDotNet installed

### Running All Benchmarks

```bash
# Navigate to the benchmarks directory
cd benchmarks

# Run all benchmarks (this will take several minutes)
go test -bench=. -benchmem -count=5
```

### Running Specific Benchmarks

```bash
# Run only client creation benchmarks
go test -bench=BenchmarkClientCreation -benchmem

# Run only message-related benchmarks
go test -bench=BenchmarkMessage -benchmem

# Run only memory diagnostics
go test -bench=BenchmarkMemoryDiagnostics -benchmem
```

### Generating Detailed Reports

```bash
# Generate a detailed report with statistics
go test -bench=. -benchmem -count=5 -cpuprofile=cpu.out -memprofile=mem.out -timeout=30m

# View the CPU profile
go tool pprof -http=:8080 cpu.out

# View the memory profile
go tool pprof -http=:8081 mem.out
```

## Benchmark Categories

### 1. Client Creation Benchmarks
Measures the overhead of creating client instances with different configurations.

- `BenchmarkClientCreation/DefaultClient` - Client with default options
- `BenchmarkClientCreation/ClientWithAPIKey` - Client with explicit API key
- `BenchmarkClientCreation/ClientWithConfig` - Client with full configuration

### 2. Message Creation Benchmarks
Measures the overhead of creating message request parameters.

- `BenchmarkMessageCreation/SimpleMessage` - Basic message with text content
- `BenchmarkMessageCreation/MessageWithSystemPrompt` - Message with system prompt
- `BenchmarkMessageCreation/MessageWithMultipleBlocks` - Message with multiple text blocks

### 3. Message Service Benchmarks
Measures the overhead of creating message service instances.

- `BenchmarkMessageServiceCreation/NewMessageService` - Message service creation

### 4. Model Operations Benchmarks
Measures performance of model-related operations.

- `BenchmarkModelOperations/ModelList` - List available models
- `BenchmarkModelOperations/ModelRetrieve` - Retrieve specific model details

### 5. Completion Operations Benchmarks
Measures performance of completion operations.

- `BenchmarkCompletionOperations/CompletionCreate` - Create completion request

### 6. Message Parameter Construction Benchmarks
Measures the overhead of creating different message parameter types.

- `BenchmarkMessageParamConstruction/UserMessageWithText` - User message with text
- `BenchmarkMessageParamConstruction/AssistantMessageWithText` - Assistant message with text
- `BenchmarkMessageParamConstruction/ToolMessage` - Tool message
- `BenchmarkMessageParamConstruction/TextBlockCreation` - Text block creation

### 7. Block Types Benchmarks
Measures performance of different block types.

- `BenchmarkBlockTypes/TextBlock` - Text block
- `BenchmarkBlockTypes/ImageBlock` - Image block
- `BenchmarkBlockTypes/DocumentBlock` - Document block

### 8. Message Validation Benchmarks
Measures the overhead of message validation.

- `BenchmarkMessageValidation/ValidateMessageParams` - Message parameter validation

### 9. Memory Diagnostics Benchmarks
Enables memory allocation tracking for critical operations.

- `BenchmarkMemoryDiagnostics/ClientCreationMemory` - Client creation memory usage
- `BenchmarkMemoryDiagnostics/MessageParamsMemory` - Message parameters memory usage
- `BenchmarkMemoryDiagnostics/BlockTypesMemory` - Block types memory usage

## Interpreting Results

Each benchmark reports:
- **ns/op**: Nanoseconds per operation (lower is better)
- **allocs/op**: Memory allocations per operation (lower is better)
- **B/op**: Bytes allocated per operation (lower is better)

### Example Output
```
BenchmarkClientCreation/DefaultClient-8          12345678   98.76 ns/op   1234 B/op   45 allocs/op
BenchmarkClientCreation/ClientWithAPIKey-8     23456789   45.67 ns/op    567 B/op   23 allocs/op
BenchmarkMessageCreation/SimpleMessage-8        34567890  123.45 ns/op   2345 B/op   67 allocs/op
```

## Adding New Benchmarks

To add new benchmarks:

1. Create a new benchmark function with the `Benchmark` prefix
2. Use `b.ResetTimer()` to reset the timer before the actual benchmark code
3. Use `b.ReportAllocs()` to track memory allocations if needed
4. Follow the existing naming conventions for benchmark groups
5. Add appropriate test data to ensure realistic measurements

## Performance Goals

The benchmarks help identify performance regressions and optimization opportunities. Target improvements:

- Client creation should be fast (< 100μs)
- Message parameter construction should be efficient (< 50μs)
- Memory allocations should be minimal (< 1KB per operation)
- No significant performance regressions between releases

## Notes

- Benchmarks are designed to run without actual API calls (no real network requests)
- Tests use mock data to avoid external dependencies
- Results may vary based on hardware and Go version
- For accurate results, run benchmarks multiple times with `-count=5`


## License

This benchmark project is part of the anthropic-sdk-go repository and is licensed under the MIT License.