package controllers

import (
	"github.com/JohnSmithX/mus/app/api"
	"github.com/JohnSmithX/mus/app/manager"
	"github.com/JohnSmithX/mus/app/db"
)



func NewAPI(store *db.Storage) {
	//do some to initialize
	//create a manager (first arg -> show debug)
	SM := manager.NewManager(true)

	servers, err := api.GetAllServersFromRedis(store)

	SM.AddServersToManager(servers)


}
