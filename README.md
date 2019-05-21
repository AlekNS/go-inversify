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
* Autowire structure
* Modules

## Examples

TODO

#### Values

```
  container := inversify.NewContainer("name of container")

  container.Bind(1).To("Hello")
  container.Build()

  container.Get(1) == "Hello"


  container.Rebind(1).To("world")
  container.Build()

  container.Get(1) == "world"
```

#### Factories, singleton and optional dependencies

```
  container.Bind(1).To("Hello")
  container.Bind(2).To(" world")

  // abstract - using of Any
  container.Bind(3).ToFactory(func (word1, word2, optDep Any) Any {
    return word1.(string) + word2.(string), nil
  }, 1, 2, inversify.Optional(4))

  // or through typed function

  container.Bind(3).ToTypedFactory(func (word1, word2 string, optDep Any) string {
    // optDep == nil
    return word1 + word2, nil
  }, 1, 2, inversify.Optional(4)).InSingletonScope()

  // resolved

  container.Build()

  container.IsBound(3) == true
  container.Get(3) == "Hello world"

```

#### Named dependencies

```
  container := inversify.NewContainer("name of container")

  container.Bind(1).To("empty")
  container.Bind(1, "dev").To("Hello")
  container.Bind(1, "prod").To("world")
  container.Build()

  container.Get(1) == "empty"
  container.Get(1, "dev") == "Hello"
  container.Get(1, "prod") == "world"
```

#### Merge

```
  mergedContainer := container1.Merge(container2)
  mergedContainer.Build()
```

#### Hierarchy

```
  baseContainer := Container()
  ...
  subContainer1 := Container()
  ...
  subContainer2 := Container()
  ...

  subContainer1.SetParent(baseContainer)
  subContainer1.Build()

  subContainer2.SetParent(baseContainer)
  subContainer2.Build()
```

#### Autowire

```
  type AppConfig struct {}

  type App struct {
    Values map[string]interface{}  `inversify:"strkey:values"`

    Config *AppConfig `inversify:""`

    TaskRepository ITaskRepository `inversify:"named:xorm"`
    Scheduler      IScheduler      `inversify:"optional"`

    Billing  BillingService // no injection
  }

  container.Bind("values").To(map[string]interface{}{ ... })
  container.Bind((*AppConfig)(nil)).To(&AppConfig{ ... })
  container.Bind((*ITaskRepository)(nil), "xorm").ToFactory(func() (Any, error) {
    ...
  })

  app := &App{}
  container.AutowireStruct(app)
```

#### Modules

```
  authModule := NewModule("auth").
      Register(func(c ContainerBinder) error {
            c.Bind()
            return nil
      }).
      UnRegister(func(c ContainerBinder) error {
            // c.Unbind()
            return nil
      })

  container.Load(authModule)
  container.Load(otherModule)
  container.Build()

  container.Unload(authModule)
  container.Unload(otherModule)
  container.Build()
```
