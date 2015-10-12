# go-profiling-hooks

This library provides single entry point to all profiling functionality available in golang 1.5. `StartProfiling` starts writing [trace](https://golang.org/cmd/trace/) and [cpu profile](https://golang.org/pkg/runtime/pprof/#StartCPUProfile) to some random directory it creates before running.
When you call `StopProfiling` it writes [heap profile](https://golang.org/pkg/runtime/pprof/#WriteHeapProfile) to the same directory as well as stopping current profiling.

## License

MIT