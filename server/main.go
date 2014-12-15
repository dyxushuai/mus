package main

import (
	"fmt"
	"net"
)

func main() {
	maxOfSS := make(map[string]chan int, 2)
	fmt.Println(maxOfSS[0])
	net.Listen("tcp", ":8044")
}
