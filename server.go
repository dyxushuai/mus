package main

import (
	"github.com/JohnSmithX/mus/manager"
	"net/http"
	"github.com/JohnSmithX/mus/controllers"
	"github.com/gohttp/app"
	"github.com/gohttp/logger"
	"github.com/goocean/methodoverride"
)

func main() {
//	m := manager.New(true)
//
//	m.AddServerAndRun("9090", "rc4", "123456", 1111111, 10)
//	http.ListenAndServe(":1234", nil)
	server := app.New()
	server.Use(logger.New())
	server.Use(methodoverride.New())

	server.Get("/api/servers", "get all")
	server.Post("/api/servers", "create new")

	server.Get("/api/servers/:id", "get :id server")
	server.Del("/api/servers/:id", "delete :id server")
	server.Put("/api/servers/:id", "update :id server")

	server.Post("/api/servers/:id/start", "start :id server")
	server.Post("/api/servers/:id/stop", "stop :id server")
	server.Post("/api/servers/:id/restart", "restart :id server")
	server.Post("/api/servers/:id/pause", "pause :id server")
	server.Post("/api/servers/:id/unpause", "pause :id server")

	server.Get("/api/servers/:id/logs", "get :id server logs")
	server.Get("/api/servers/:id/flow", "get :id server flow")




}
