# Go-Inversify - Dependency Injection Tool for Go

[![Build Status](https://travis-ci.org/AlekNS/go-inversify.svg?branch=master)](https://travis-ci.org/AlekNS/go-inversify)
[![Go Report Card](https://goreportcard.com/badge/github.com/AlekNS/go-inversify)](https://goreportcard.com/report/github.com/AlekNS/go-inversify)

Yet another dependency injection system for Go inspired by [InversifyJS](https://github.com/inversify/InversifyJS).

## Installation

#### Go get

```
go get github.com/alekns/go-inversify
```

## Examples

#### Values

```
  container := inversify.Container()

  container.Bind(1).To("Hello")
  container.Bind(2).To(" world")
  container.Bind(3).ToFactory(func (word1, word2, optDep Any) Any {
    return word1.(string) + word2.(string), nil
  }, 1, 2, inversify.Optional(4))

  // or

  container.Bind(3).ToTypedFactory(func (word1, word2 string, optDep Any) string {
    // optDep == nil
    return word1 + word2, nil
  }, 1, 2, inversify.Optional(4))

  container.IsBound(3) == true
  container.Get(3).(string) == "Hello world"
```
