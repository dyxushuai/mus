package app


import (
	"github.com/JohnSmithX/mus/app/controllers"
	"github.com/gohttp/app"
	"github.com/gohttp/logger"
	"github.com/goocean/methodoverride"
	"github.com/JohnSmithX/mus/app/utils"

)

func Serve(redisHost, redisPWD string) {
	server := app.New()
	server.Use(logger.New())
	server.Use(methodoverride.New())

	api := controllers.New(redisHost, redisPWD)
	//	server.Get("/api/servers", "get all")
	//	server.Post("/api/servers", "create new")

	server.Get("/api/servers/:id", utils.JsonView(api.NewServerAPI().Show))
	//	server.Del("/api/servers/:id", "delete :id server")
	//	server.Put("/api/servers/:id", "update :id server")
	//
	//	server.Post("/api/servers/:id/start", "start :id server")
	//	server.Post("/api/servers/:id/stop", "stop :id server")
	//	server.Post("/api/servers/:id/restart", "restart :id server")
	//
	//	server.Get("/api/servers/:id/logs", "get :id server logs")
	//	server.Get("/api/servers/:id/flow", "get :id server flow")

	server.Listen(":7888")

}
