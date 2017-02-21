Golang Logging Library
=======================

The colorful and simple logging library for Golang.

## What is this

Opinionated logging library that able to output to `io.Reader` with file descriptors (`os.Stdout`, `os.Stderr`,
regular file, etc.) with automatic terminal color support.

The log serverity and behavior can be described as follows:

| Severity | Description | Caller Info |
|:--------:|:------------|:-----------:|
| Fatal | Unrecoverable error and automatically exit after logging | Yes |
| Error | Recoverable error but need attention | Yes |
| Warn | Minor error and does not output the caller info | No |
| Info | Informational message | No |
| Debug | Debug message, only shown when debug enabled | Yes |
| Trace | Trace message, only shown when debug enabled | No |

## Getting started

Add the `go-log` package using

```
go get github.com/withmandala/go-log
```

And import it to your package by

```go
import (
    "github.com/withmandala/go-log"
)
```

Use the `go-log` package with

```go
logger := log.New(os.Stderr)
logger.Info("Hi, this is your logger")
```

## Color support

The library will try to automatically detect the `io.Reader` file descriptor when calling `log.New()` for color
support. But, if you insist to use or not to use color, you can add `.WithColor()` or `.WithoutColor()` respectively.

```go
// With color
logger := log.New(os.Stderr).WithColor()

// Without color
logger := log.New(os.Stderr).WithoutColor()
```

## Debug output

The log library will suppress the `.Debug()` and `.Trace()` output by default. To enable or disable the debug output,
call `(Logger).WithDebug()` or `(Logger).WithoutDebug()` respectively.

```go
// Enable debugging
logger := log.New(os.Stderr).WithDebug()
// Print debug output
logger.Debug("Test debug output")
// Disable debug output
logger.WithoutDebug()
logger.Debug("Test debug output") // This message will not be printed
```

## Be Quiet

If somehow the log is annoying to you, just shush it by calling `(Logger).Quiet()` and **ALL** log output will be
disappear, although `.Fatal()` will silently quit the program with error. To re-enable the log output use
`(Logger).NoQuiet()`.
