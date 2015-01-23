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

	controllers.NewAPI(redisHost, redisPWD)

	serverAPI := controllers.NewServerAPI()

	server.Get("/api/servers", utils.JsonView(serverAPI.Index))
	server.Post("/api/servers", utils.JsonView(serverAPI.Create))

	server.Get("/api/servers/:id", utils.JsonView(serverAPI.Show))
	server.Del("/api/servers/:id", utils.JsonView(serverAPI.Destroy))
	server.Put("/api/servers/:id", utils.JsonView(serverAPI.Update))



	//	server.Post("/api/servers/:id/start", "start :id server")
	//	server.Post("/api/servers/:id/stop", "stop :id server")
	//	server.Post("/api/servers/:id/restart", "restart :id server")
	//
	//	server.Get("/api/servers/:id/logs", "get :id server logs")
	//	server.Get("/api/servers/:id/flow", "get :id server flow")

	server.Listen(":7888")

}
