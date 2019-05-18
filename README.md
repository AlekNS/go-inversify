# Go-Inversify - Dependency Injection Tool for Go

[![Build Status](https://travis-ci.org/AlekNS/go-inversify.svg?branch=master)](https://travis-ci.org/AlekNS/go-inversify)
[![Go Report Card](https://goreportcard.com/badge/github.com/AlekNS/go-inversify)](https://goreportcard.com/report/github.com/AlekNS/go-inversify)

Yet another dependency injection system for Go inspired by [InversifyJS](https://github.com/inversify/InversifyJS).

## Installation

#### Go get

```
go get github.com/alekns/go-inversify
```

#### Features

* Bind (and rebind) any value and types to any values
* Abstract factory (in terms of interface{} - fast), typed factory (slow) with normal types
* Singletons (lazy)
* Optional dependencies (resolved as nil)
* Named dependencies (multi-bindings on single value or type)
* Checking on cycle dependencies (panic on Build)
* Containers merging
* Containers hierarchy
* Modules

## Examples

TODO

#### Values

```
  container := inversify.Container()

  container.Bind(1).To("Hello")
  container.Build()

  container.Get(1).(string) == "Hello"
```

#### Factories, singleton and optional dependencies

```
  container.Bind(1).To("Hello")
  container.Bind(2).To(" world")

  container.Bind(3).ToFactory(func (word1, word2, optDep Any) Any {
    return word1.(string) + word2.(string), nil
  }, 1, 2, inversify.Optional(4))

  // or

  container.Bind(3).ToTypedFactory(func (word1, word2 string, optDep Any) string {
    // optDep == nil
    return word1 + word2, nil
  }, 1, 2, inversify.Optional(4)).InSingletonScope()

  // resolved

  container.Build()

  container.IsBound(3) == true
  container.Get(3).(string) == "Hello world"

```

#### Named dependencies

```
  container := inversify.Container()

  container.Bind(1).To("empty")
  container.Bind(1, "dev").To("Hello")
  container.Bind(1, "prod").To("world")
  container.Build()

  container.Get(1).(string) == "empty"
  container.Get(1, "dev").(string) == "Hello"
  container.Get(1, "prod").(string) == "world"
```

#### Merge

```
TODO example
```

#### Hierarchy

```
TODO example
```

#### Modules

```
TODO example
```
