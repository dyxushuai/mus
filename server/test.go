package main

import (
	"fmt"

)
func foo() (e int) {

	a , e, w := 4, 6, 7

	fmt.Println(a)
	return
}

func main() {
	fmt.Println(foo())
}

