# goprof

This library provides single entry point to all profiling functionality available in golang 1.5. `StartProfiling` starts writing [trace](https://golang.org/cmd/trace/) and [cpu profile](https://golang.org/pkg/runtime/pprof/#StartCPUProfile) to some random directory it creates before running.
When you call `StopProfiling` it writes [heap profile](https://golang.org/pkg/runtime/pprof/#WriteHeapProfile) to the same directory as well as stopping current profiling.

## Logging

By default, the library writes logs about start/stop profiling and errors using standard go logger. You can provide
your own log function in order to make it fit your logging or shut up the logging at all.

We don't use much log levels since all the messages have quite the same level.

## License

MIT