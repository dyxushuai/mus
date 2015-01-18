package controllers

import (
	"github.com/JohnSmithX/mus/app/models"
	"github.com/JohnSmithX/mus/app/manager"
	"github.com/JohnSmithX/mus/app/db"
	"github.com/JohnSmithX/mus/app/utils"
	"fmt"

)

type API struct {
	SM *manager.Manager
}

func New(redisHost, redisPWD string) (api *API) {

	store := db.NewStorage(redisHost, redisPWD)
	fmt.Println(redisHost, redisPWD)
	//do some to initialize
	//create a manager (first arg -> show debug)
	api = &API{}
	api.SM = manager.NewManager()

	servers, err := models.GetAllServersFromRedis(store)
	if err != nil {
		utils.Debug(err)
	}

	for _, server := range servers {
		api.SM.AddServerToManager(server)
	}


	return
}


func (self *API) NewServerAPI() *ServerAPI {
	return &ServerAPI{
		SM: self.SM,
	}
}
