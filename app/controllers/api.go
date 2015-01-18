package controllers

import (
	"github.com/JohnSmithX/mus/app/models"
	"github.com/JohnSmithX/mus/app/manager"
	"github.com/JohnSmithX/mus/app/db"
	"github.com/JohnSmithX/mus/app/utils"
)



func New(store *db.Storage) {
	//do some to initialize
	//create a manager (first arg -> show debug)
	SM := manager.NewManager()

	servers, err := models.GetAllServersFromRedis(store)
	if err != nil {
		utils.Debug(err)
	}
	SM.AddServersToManager(servers)


}
