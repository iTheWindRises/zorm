package main

import (
	"fmt"
	"reflect"
)

type User struct {
	Name string `zorm:"PRIMARY KEY"`
	Age int
}

func main() {

	fmt.Println(reflect.Indirect(reflect.ValueOf(&User{})).Type().String())
	fmt.Println(reflect.Indirect(reflect.ValueOf(User{})).Type().String())

}
