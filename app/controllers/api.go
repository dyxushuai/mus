package controllers

import (
	"github.com/JohnSmithX/mus/app/models"
	"github.com/JohnSmithX/mus/app/manager"

	"github.com/JohnSmithX/mus/app/utils"
)
var (
	SM *manager.Manager
)




func NewAPI(redisHost, redisPWD string) {

	//do some to initialize
	//create a manager (first arg -> show debug)

	SM = manager.NewManager()

	models.InitDb(redisHost, redisPWD)

	servers, err := models.GetAllServersFromRedis()
	if err != nil {
		utils.Debug(err)
	}


	for _, server := range servers {
		SM.AddServerToManager(server)
	}


	return
}


func NewServerAPI() *ServerAPI {
	return &ServerAPI{
	}
}
