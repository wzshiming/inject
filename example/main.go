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
