Go-Log
=======

The colorful and simple logging library for Golang

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
