package main

import (
	"github.com/JohnSmithX/mus/manager"
	"github.com/JohnSmithX/mus/models"
//	"github.com/gohttp/app"
//	"github.com/gohttp/logger"
//	"github.com/goocean/methodoverride"
)

func main() {
	m := manager.New()

	models.NewStorage("")
	s := models.Server{
		Port: "9090",
		Password: "123456",
		Method:   "rc4",
	}
	m.AddServerAndRun(&s)
	go m.LOG()
	m.DEBUG()
//	server := app.New()
//	server.Use(logger.New())
//	server.Use(methodoverride.New())

}
