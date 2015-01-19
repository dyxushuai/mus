package controllers

import (
	"github.com/JohnSmithX/mus/app/models"
	"github.com/JohnSmithX/mus/app/manager"
	"github.com/JohnSmithX/mus/app/db"
	"github.com/JohnSmithX/mus/app/utils"
	"fmt"

)
var (
	Store db.IStorage
	SM *manager.Manager
)




func Initialize(redisHost, redisPWD string) {

	store := db.NewStorage(redisHost, redisPWD)
	fmt.Println(redisHost, redisPWD)
	//do some to initialize
	//create a manager (first arg -> show debug)

	SM = manager.NewManager()
	Store = store

	servers, err := models.GetAllServersFromRedis(store)
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
