package app


import (
	"github.com/JohnSmithX/mus/app/controllers"
	"github.com/gohttp/app"
	"github.com/gohttp/logger"
	"github.com/goocean/methodoverride"
	"github.com/JohnSmithX/mus/app/db"
	"github.com/JohnSmithX/mus/app/middlewares"

)

var (
	rdPool *db.Storage
)

func NewRedisPool(serverHost, serverPassword string) *db.Storage {
	return db.NewStorage(serverHost, serverPassword)
}


func Serve(redisHost, redisPWD string) {
	rdPool = NewRedisPool(redisHost, redisPWD)

	server := app.New()
	server.Use(logger.New())
	server.Use(middlewares.Auth(rdPool))
	server.Use(methodoverride.New())


	controllers.NewAPI(rdPool)

	serverAPI := controllers.NewServerAPI()
	serverActions := controllers.NewAction()

	server.Get("/api/servers", controllers.JsonView(serverAPI.Index))
	server.Post("/api/servers", controllers.JsonView(serverAPI.Create))

	server.Get("/api/servers/:id", controllers.JsonView(serverAPI.Show))
	server.Del("/api/servers/:id", controllers.JsonView(serverAPI.Destroy))
	server.Put("/api/servers/:id", controllers.JsonView(serverAPI.Update))



	server.Post("/api/servers/:id/start", controllers.JsonView(serverActions.Start))
	server.Post("/api/servers/:id/stop", controllers.JsonView(serverActions.Stop))
	server.Post("/api/servers/:id/restart", controllers.JsonView(serverActions.Restart))
	//
	//	server.Get("/api/servers/:id/logs", "get :id server logs")
	//	server.Get("/api/servers/:id/flow", "get :id server flow")

	server.Listen(":7888")

}
