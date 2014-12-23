package main

import (
	"github.com/JohnSmithX/mus/server/manager"
//	"fmt"
)

func main() {
	m := manager.CreateManager()
	m.AddServerAndRun("9090")
//	m.AddServerAndRun("9090")
	m.AddServerAndRun("8080")
	m.StopServer("8080")
	m.StartServer("8080")
	m.ServerList()
	go m.LOG()
	m.DEBUG()

}
