# YALA - Yet Another Logging Abstraction for Go

[![Build](https://github.com/elgopher/yala/actions/workflows/build.yml/badge.svg)](https://github.com/elgopher/yala/actions/workflows/build.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/elgopher/yala.svg)](https://pkg.go.dev/github.com/elgopher/yala)
[![Go Report Card](https://goreportcard.com/badge/github.com/elgopher/yala)](https://goreportcard.com/report/github.com/elgopher/yala)
<img src="logo.png" align="right" width="30%">

Tiny **structured logging** abstraction or facade for various logging libraries, allowing the end user to plug in the desired logging library in `main.go`.

## Supported logging libraries (via adapters)

[logrus](adapter/logrusadapter), [zap](adapter/zapadapter), [zerolog](adapter/zerologadapter), [glog](adapter/glogadapter), [log15](adapter/log15adapter) and [standard fmt and log packages](adapter/printer)

## When to use?

* If you are a package/module/library author
* And you want to participate in the end user logging system (log messages using the logger provided by the end user)
* You don't want to add dependency to any specific logging library to your code
* You don't want to manually inject logger to every possible place where you want to log something (such as function, struct etc.)
* If you need a nice and elegant API with a bunch of useful functions, but at the same time you don't want your end users to spend hours on writing their own logging adapter.

## Installation

```shell
# Add yala to your Go module:
go get github.com/elgopher/yala        
```

Please note that at least Go `1.17` is required.

## How to use

### Choose logger - global or local?

Global logger can be accessed from everywhere in your library and can be reconfigured anytime. Local logger is a logger
initialized only once and used locally, for example inside the function.

### Use global logger

```go
package lib // this is your package, part of module/library etc.

import (
	"context"
	"errors"

	"github.com/elgopher/yala/logger"
)

// define global logger, no need to initialize it (by default nothing is logged)
var log logger.Global

// Provide a public function for setting adapter. It will be called in main.go
func SetLogAdapter(adapter logger.Adapter) {
	log.SetAdapter(adapter)
}

func Function(ctx context.Context) {
	log.Debug(ctx, "Debug message")
	log.With(ctx, "field_name", "value").Info("Message with field")
	log.WithError(ctx, errors.New("some")).Error("Message with error")
}
```

#### Specify adapter - a real logger implementation.

```go
package main

import (
	"context"

	"github.com/elgopher/yala/adapter/printer"
	"lib"
)

// End user decides what library to plug in.
func main() {
	adapter := printer.StdoutAdapter() // will use fmt.Println
	lib.SetLogAdapter(adapter)

	ctx := context.Background()
	lib.Function(ctx)
}
```

### Why context.Context is a parameter?

`context.Context` can very useful in transiting request-scoped tags or even entire logger. A `logger.Adapter` implementation might use them
making possible to log messages instrumented with tags. Thanks to that your library can trully participate in the incoming request. 

### Use local logger

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


// end user code
adapter := printer.StdoutAdapter()
lib := NewLibrary(adapter)
```

### How to use existing adapters

* [Logrus](adapter/logrusadapter/_example/main.go)
* [fmt.Println and standard log package](adapter/printer/_example/main.go)
* [Zap](adapter/zapadapter/_example/main.go)
* [Zerolog](adapter/zerologadapter/_example/main.go)
* [glog](adapter/glogadapter/_example/main.go)
* [Log15](adapter/log15adapter/_example/main.go)

### Writing your own adapter

Just implement `logger.Adapter` interface:

```go
type MyAdapter struct{}

func (MyAdapter) Log(context.Context, logger.Entry) {
    // here you can do whatever you want with the log entry 
}
```

### Difference between Logger and Adapter

* Logger is used by package/module/library author
* Adapter is an interface to be implemented by adapters. They use real logging libraries under the hood.
* So, why two abstractions? Simply because the smaller the Adapter interface, the easier it is to implement it. On the other hand, from library perspective, more methods means API which is easier to use. 

### Advanced recipes

* [Filter out messages starting with given prefix](logger/_examples/filter/main.go)
* [Add field to each message taken from context.Context](logger/_examples/tags/main.go)
* [Zap logger passed over context.Context](logger/_examples/contextlogger/main.go)

### Why just don't create my own abstraction instead of using yala?

Yes, you can also create your own. Very often it is just an interface with a single method, like this:

```go
type ImaginaryLogger interface {
    Log(context.Context, Entry)
}
```

But there are limitations for such solution:

* such interface alone is not very easy to use in your package/module/library
* someone who is using your package is supposed to write his own adapter of this interface (or you can provide adapters which
  of course takes your valuable time)
* it is not obvious how logging API should look like

### But yala is just another API. Why it is unique?

* yala is designed for the ease of use. And by that I mean ease of use for everyone - developer logging messages, developer writing adapter and end user configuring the adapter
* yala is using `context.Context` in each method call, making possible to use sophisticated request-scoped logging

### YALA limitations

* even though your package will be independent of any specific logging implementation, you still have to import 
  `github.com/elgopher/yala/logger`. This package is relatively small though, compared to real logging libraries
  (about ~200 lines of production code) and **it does not import any external libraries**.
* yala is not optimized for **extreme** performance, because this would hurt the developer experience and readability of the created code. Any intermediary API ads overhead - global synchronized variables, wrapper code and even polymorphism slow down the execution a bit. The overhead varies, but it is usually a matter of tens of nanoseconds per call. 
