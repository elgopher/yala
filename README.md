[![Go](https://github.com/jacekolszak/yala/actions/workflows/go.yml/badge.svg)](https://github.com/jacekolszak/yala/actions/workflows/go.yml)
![logo](yala.png)
# YALA - Yet Another Logging Abstraction for Go

Tiny structured logging abstraction with adapters for most popular logging Go libraries and easy way to roll your own.

## Supported logging implementations

[logrus](adapter/logrusadapter), [zap](adapter/zapadapter), [zerolog](adapter/zerologadapter), [glog](adapter/glogadapter), [log15](adapter/log15adapter) and [standard fmt and log packages](adapter/printer)

## When to use?

* If you are a module/package/library author
* And you want to participate in a caller logging system (log messages using the logger provided by the consumer)
* And you don't want to add dependency to any specific logging library to your code
* And you don't want to manually inject logger to every possible place where you want to log something (function, struct etc.)
* If you need nice and elegant API with a bunch of useful functions, but at the same time you don't want your clients spend hours on writing their own logging adapter.

## Installation

```shell
# Add yala to your Go module:
go get github.com/jacekolszak/yala        
```

## How to use

### Choose logger - global or local?

Global logger can be accessed from everywhere in your library and can be reconfigured anytime. Local logger is a logger
initialized only once and used locally, for example inside the function.

### Use global logger

```go
package lib // this is your package, part of module/library etc.

import (
	"context"

	"github.com/jacekolszak/yala/logger"
)

var Logger logger.Global // define global logger, no need to initialize (by default nothing is logged)

func Function(ctx context.Context) {
	Logger.Debug(ctx, "Debug message")
	Logger.With(ctx, "field_name", "value").Info("Message with field")
	Logger.WithError(ctx, err).Error("Message with error")
}
```

#### Specify adapter - a real logger implementation.

```go
package main

import (
  "context"

  "github.com/jacekolszak/yala/adapter/printer"
  "lib"
)

func main() {
  adapter := printer.StdoutAdapter() // will use fmt.Println
  lib.Logger.SetAdapter(adapter)     // set the adapter

  ctx := context.Background()
  lib.Function(ctx)
}
```

### Why context.Context is a parameter?

`context.Context` can very useful in transiting request-scoped tags or even entire logger. A `logger.Adapter` implementation might use them
making possible to log messages instrumented with tags. Thanks to that your library can trully participate in the incoming request. 

### Why global state?

Logging is a special kind of dependency. It is used all over the place. Adding it as an explicit dependency to every
function, struct etc. can be cumbersome. Still though, you have an option to use **local** logger by injecting
logger.Adapter into your library:

```go
// your library code:
func NewLibrary(adapter logger.Adapter) YourLib {
    // create a new local logger which provides similar API to the global logger
    localLogger := logger.Local(adapter)         
    return YourLib{localLogger: localLogger}
}

type YourLib struct {
	localLogger logger.LocalLogger
}

func (l YourLib) Method(ctx context.Context) {
    l.localLogger.Debug(ctx, "message from local logger")
}

// client code
adapter := printer.StdoutAdapter()
lib := NewLibrary(adapter)
```

### More examples

* [Logrus](adapter/logrusadapter/_example/main.go)
* [fmt.Println and standard log package](adapter/printer/_example/main.go)
* [Zap](adapter/zapadapter/_example/main.go)
* [Zerolog](adapter/zerologadapter/_example/main.go)
* [glog](adapter/glogadapter/_example/main.go)
* [Log15](adapter/log15adapter/_example/main.go)
* [Zap logger passed over context.Context](adapter/contextadapter/_example/main.go)

### Writing your own adapter

Just implement `logger.Adapter` interface:

```go
type MyAdapter struct{}

func (MyAdapter) Log(context.Context, logger.Entry) {
    // here you can do whatever you want with the log entry 
}
```

### Difference between logger.Logger and logger.Adapter

* logger.Logger is a struct for logging messages (optionally with fields and error). It is used by packages in your module.
* logger.Adapter is an abstraction which should be implemented by adapters. Some adapters are already implemented (such as logrusadapter) and new adapter can be easily implemented too.
* so, why two abstractions? Simply because the smaller the Adapter interface, the easier it is to implement it. On the other hand, from library perspective, more methods means API which is easier to use. 

### Why just don't create my own abstraction instead of using yala?

Yes, you can also create your own. Very often it just an interface with a single method, like this:

```go
type ImaginaryLogger interface {
    Log(context.Context, Entry)
}
```

But there are limitations for such solution:

* such interface alone is not very easy to use in your module/library/package
* someone who is using your module is supposed to write his own adapter of this interface (or you can provide adapters which
  of course takes your valuable time)
* it is not obvious how logging API should look like

### YALA limitations

* even though your module will be independent of any specific logging implementation, you still have to import 
  `github.com/jacekolszak/yala/logger`. This package is relatively small though, compared to real logging libraries
  (about ~200 lines of production code) and **it does not import any external libraries**.
* yala is optimized for the ease of use (both for the developer who logs messages and for the developer writing
  adapter). It is not optimized for *extreme* performance, because this would hurt the user experience and readability of the
  created code. Any intermediary API ads overhead - global synchronized variables, wrapper code and even polymorphism slow down 
  the execution a bit. The overhead varies, but it is usually a matter of tens of nanoseconds per call. 
