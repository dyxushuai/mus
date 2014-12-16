package main

import (
"fmt"


)
type A struct {
	ax, ay int
}

const (
	a int = iota
	b
	c
)
type B struct {
	A
	bx, by float64
}
func main() {
	fmt.Println(b == 2)

}
