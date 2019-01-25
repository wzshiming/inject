# Inject

[![Build Status](https://travis-ci.org/wzshiming/inject.svg?branch=master)](https://travis-ci.org/wzshiming/inject)
[![Go Report Card](https://goreportcard.com/badge/github.com/wzshiming/inject)](https://goreportcard.com/report/github.com/wzshiming/inject)
[![GoDoc](https://godoc.org/github.com/wzshiming/inject?status.svg)](https://godoc.org/github.com/wzshiming/inject)
[![GitHub license](https://img.shields.io/github/license/wzshiming/inject.svg)](https://github.com/wzshiming/inject/blob/master/LICENSE)
[![cover.run](https://cover.run/go/github.com/wzshiming/inject.svg?style=flat&tag=golang-1.10)](https://cover.run/go?tag=golang-1.10&repo=github.com%2Fwzshiming%2Finject)

- [English](https://github.com/wzshiming/inject/blob/master/README.md)
- [简体中文](https://github.com/wzshiming/inject/blob/master/README_cn.md)

## Example

``` golang

package main

import (
	"fmt"
	"reflect"

	"github.com/wzshiming/inject"
)

func main() {

	inj := inject.NewInjector(nil)
	inj.Map(reflect.ValueOf(10))
	inj.Map(reflect.ValueOf("Hello world"))

	inj.Call(reflect.ValueOf(func(i int, s string) {
		fmt.Println(s, i)
		// Hello world 10
	}))

	t := struct {
		I int
		S string
	}{}
	inj.InjectStruct(reflect.ValueOf(&t))
	fmt.Println(t)
	// {10 Hello world}
}

```

## License

Pouch is licensed under the MIT License. See [LICENSE](https://github.com/wzshiming/inject/blob/master/LICENSE) for the full license text.
