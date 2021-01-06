# Bleep

[![GoDoc](https://godoc.org/github.com/sinhashubham95/bleep?status.svg)](https://pkg.go.dev/github.com/sinhashubham95/bleep)
[![Release](https://img.shields.io/github/v/release/sinhashubham95/bleep?sort=semver)](https://github.com/sinhashubham95/bleep/releases)
[![Report](https://goreportcard.com/badge/github.com/sinhashubham95/bleep)](https://goreportcard.com/report/github.com/sinhashubham95/bleep)
[![Coverage Status](https://coveralls.io/repos/github/sinhashubham95/bleep/badge.svg?branch=master)](https://coveralls.io/github/sinhashubham95/bleep?branch=master)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#utilities)

`Bleep` is used to peform actions on OS signals. It is highly extensible and goroutine safe. It is possible to add any number of actions and all of them are guaranteed to be performed simultaneously on the OS signals that `Bleep` will be listening for.

## Installation

```shell
go get github.com/sinhashubham95/bleep
```

## How to Use

The `Bleep` package allows you to create a new instance of the handler and also has a default handler in place that can be used directly.

Creating separate Bleep instances can be useful, when you want to perform different set of actions for different set of OS signals.

### Create a New OS Signal Handler

This is used to create a new handler for performing actions on OS Signals.

```go
import (
  "os"
  "github.com/sinhashubham95/bleep"
)

func New() {
  handler := bleep.New()
  // now this handler can be used to add or remove actions and listen to the OS signals
}
```

### Add an Action

This is used to add an action to be executed on the OS signal listening for.

```go
import (
  "os"
  "github.com/sinhashubham95/bleep"
)

fun Add() {
  key := bleep.Add(func (s os.Signal) {
    // do something
  })
  // this key is the unique identifier for your added action
}
```

### Remove an Action

This is used to remove an action added to Bleep.

```go
import (
  "github.com/sinhashubham95/bleep"
)

func Remove() {
  action := bleep.Remove("some-key")  // this key should be the same as the one returned during adding the action
  // the returned action is the one that was added using this key
}
```

### Listen

This is used to listen for the OS signals. Note that this will wait for the signal in the go routine in which this is called.

```go
import (
  "syscall"
  "github.com/sinhashubham95/bleep"
)

func Listen() {
  bleep.Listen(syscall.SIGINT, syscall.SIGTERM)
}
```