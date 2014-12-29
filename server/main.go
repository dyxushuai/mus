package main

import (
	"github.com/JohnSmithX/mus/server/manager"
//	"fmt"
)

func main() {
	m := manager.CreateManager()
	m.AddServerAndRun("9090")
	go m.LOG()
	m.DEBUG()

}
