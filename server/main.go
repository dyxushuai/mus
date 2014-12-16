package main

import (
//"fmt"
//	"github.com/JohnSmithX/mus/server/manager"


	"log"
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
	a := make(chan int)
	go func() {
		for {
			select{
			case com := <- a:
				switch com {
				case 1:
//					continue
				}
			}
			log.Println("here")
		}
	}()

	for {
		a <- 1
	}


}
