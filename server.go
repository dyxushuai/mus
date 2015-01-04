package main

import (
	"github.com/JohnSmithX/mus/manager"
	"net/http"
//	"github.com/JohnSmithX/mus/models"
//	"github.com/gohttp/app"
//	"github.com/gohttp/logger"
//	"github.com/goocean/methodoverride"
)

func main() {
	m := manager.New(true)

	m.AddServerAndRun("9090", "rc4", "123456", 1111111, 10)
	http.ListenAndServe(":1234", nil)
//	server := app.New()
//	server.Use(logger.New())
//	server.Use(methodoverride.New())

}
