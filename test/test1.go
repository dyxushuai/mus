package main

import (
	"fmt"
)

type A struct {
	name string
}



type B A

func main() {
	var a = A{name: "xus"}
	var b B = a
	fmt.Println(b)
}
