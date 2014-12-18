package main

import (
	"github.com/JohnSmithX/mus/server/manager"
)

func main() {
	m := manager.CreateManager()
	m.AddServerAndRun("9090")
	m.DEBUG()
}
