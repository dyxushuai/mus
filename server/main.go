package main

import (
	"net"

)
type A struct {
	ax, ay int
}

type B struct {
	A
	bx, by float64
}
func main() {
	conn, _ := net.Dial("tcp", "127.0.0.1:8080")
	conn.Close()

}
