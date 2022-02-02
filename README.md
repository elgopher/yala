# YALA - Yet Another Logging Abstraction for Go

[![Build](https://github.com/elgopher/yala/actions/workflows/build.yml/badge.svg)](https://github.com/elgopher/yala/actions/workflows/build.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/elgopher/yala.svg)](https://pkg.go.dev/github.com/elgopher/yala)
[![Go Report Card](https://goreportcard.com/badge/github.com/elgopher/yala)](https://goreportcard.com/report/github.com/elgopher/yala)
[![codecov](https://codecov.io/gh/elgopher/yala/branch/master/graph/badge.svg)](https://codecov.io/gh/elgopher/yala)
[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#active)
<img src="docs/logo.png" align="right" width="30%">

Tiny **structured logging** abstraction or facade for various logging libraries, allowing the end user to plug in the desired logging library in `main.go`.

## Supported logging libraries (via adapters)

[logrus](adapter/logrusadapter), [zap](adapter/zapadapter), [zerolog](adapter/zerologadapter), [glog](adapter/glogadapter), [log15](adapter/log15adapter), [standard log](adapter/logadapter) and [console](adapter/console)

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

### Choose logger - global or normal?

Global logger can be accessed from everywhere in your package and can be reconfigured anytime. Normal logger is an
immutable logger, initialized only once. 

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
func SetLoggerAdapter(adapter logger.Adapter) {
	log.SetAdapter(adapter)
}

func Function(ctx context.Context) {
	log.Debug(ctx, "Debug message")
	
	log.With("field_name", "value").
		Info(ctx, "Message with field")
	
	log.WithError(errors.New("some")).
		Error(ctx, "Message with error")
}
```

#### Specify adapter - a real logger implementation.

```go
package main

import (
	"context"

	"github.com/elgopher/yala/adapter/console"
	"lib"
)

// End user decides what library to plug in.
func main() {
	adapter := console.StdoutAdapter() // will print messages to console
	lib.SetLoggerAdapter(adapter)

	ctx := context.Background()
	lib.Function(ctx)
}
```

### Why context.Context is a parameter?

`context.Context` can very useful in transiting request-scoped tags or even entire logger. A `logger.Adapter` implementation might use them
making possible to log messages instrumented with tags. Thanks to that your library can trully participate in the incoming request. 

### Use normal logger

Logging is a special kind of dependency. It is used all over the place. Adding it as an explicit dependency to every
function, struct etc. can be cumbersome. Still though, you have an option to use **normal** logger by injecting
logger.Adapter into your library:

```go
// your library code:
func NewLibrary(adapter logger.Adapter) YourLib {
	// create a new normal logger which provides similar API to the global logger
	l := logger.WithAdapter(adapter)     
	return YourLib{log: l}
}

type YourLib struct {
	log logger.Logger
}

func (l YourLib) Method(ctx context.Context) {
	l.log.Debug(ctx, "message from normal logger")
}


// end user code
adapter := console.StdoutAdapter()
lib := NewLibrary(adapter)
```

### How to use existing adapters

* [Logrus](adapter/logrusadapter/_example/main.go)
* [standard log package](adapter/logadapter/_example/main.go)
* [print logs to console using simplified adapter](adapter/console/_example/main.go)
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
* Here is the architecture from the package perspective:

<img src="docs/architecture.svg" width="100%">


### More examples

* [How to reuse logger](logger/_examples/reuse/main.go)

### Advanced recipes

* [Filter out messages starting with given prefix](logger/_examples/filter/main.go)
* [Filter messages by level](logger/_examples/levelfilter/main.go)
* [Add field to each message taken from context.Context](logger/_examples/tags/main.go)
* [Report caller information in each message](logger/_examples/caller/main.go)
* [Zap logger passed over context.Context](logger/_examples/contextlogger/main.go)

## Why just don't create my own abstraction instead of using yala?

Yes, you can also create your own. Very often it is just an interface with a single method, like this:

```go
type ImaginaryLogger interface {
    Log(context.Context, Entry)
}
```

But there are limitations for such solution:

* such interface alone is not very easy to use in your package/module/library. You just have to write way too much boilerplate code.
* someone who is using your package is supposed to write implementation of this interface (or you can provide prebuilt implementation for various logging libraries). In both cases this cost time and effort.
* it is not obvious how logging API should look like. Someone would argue that is better to have a much more complicated interface like this:
```go
type AnotherImaginaryLogger interface {
	With(field string, value interface{}) AnotherImaginaryLogger
	WithError(err error) AnotherImaginaryLogger
	Info(context.Context, string)
	Debug(context.Context, string)
	Warn(context.Context, string)
	Error(context.Context, string)
}
```
Unfortunately such interface is much harder to implement, than interface with a single method.

## But yala is just another API. Why is it unique?

* yala is designed for the ease of use. And by that I mean ease of use for everyone - developer logging messages, developer writing adapter and end user configuring the adapter:
  * two types of concurrency-safe loggers
  * easy to implement, one-method adapter interface
  * full control over what is logged, and how
* yala is using `context.Context` in each method call, making possible to use sophisticated request-scoped logging

## YALA limitations

* even though your package will be independent of any specific logging implementation, you still have to import 
  `github.com/elgopher/yala/logger`. This package is relatively small though, compared to real logging libraries
  (about ~200 lines of production code) and **it does not import any external libraries**.
* yala is not optimized for **extreme** high performance, because this would hurt the developer experience and readability of the created code. Any intermediary API ads overhead - global synchronized variables, wrapper code and even polymorphism slow down the execution a bit. The overhead varies, but it is usually a matter of tens of nanoseconds per call. 
