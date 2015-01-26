package controllers

import (
	"github.com/JohnSmithX/mus/app/models"
	"github.com/JohnSmithX/mus/app/manager"
	"github.com/JohnSmithX/mus/app/utils"
	"net/http"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"github.com/JohnSmithX/mus/app/db"
)
var (
	SM *manager.Manager
)




func NewAPI(rdPool *db.Storage) {

	//do some to initialize
	//create a manager (first arg -> show debug)

	SM = manager.NewManager()

	models.InitDb(rdPool)
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

func NewAction() *ServerActionsAPI {
	return &ServerActionsAPI{

	}
}

func JsonView(fn func(w http.ResponseWriter, r *http.Request) (string, error)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		data, err := fn(w, r)
		if err != nil {

			if v, ok := err.(errors.DropboxError); ok {
				e, _ := models.NewErr(v.GetMessage()).JSON()
				data = string(e)
			} else {
				e, _ := models.NewErr(err.Error()).JSON()
				data = string(e)
			}

		}
		w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=%s", "application/json", "utf-8"))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		fmt.Fprintf(w, data)

	}
}
