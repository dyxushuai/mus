package controllers

import (
	"net/http"
	"github.com/JohnSmithX/mus/app/manager"
	"github.com/JohnSmithX/mus/app/utils"
)

type ServerAPI struct {
	SM *manager.Manager
}




//Get("/api/servers", "get all")
func (self *ServerAPI) Index(w http.ResponseWriter, r *http.Request) (json string) {
	return

}

//Post("/api/servers", "create new")
func (self *ServerAPI) Create(w http.ResponseWriter, r *http.Request) (json string) {
	return
}

//Get("/api/servers/:id", "get :id server")
func (self *ServerAPI) Show(w http.ResponseWriter, r *http.Request) (json string) {
	params := r.URL.Query()
	id := params.Get(":id")


	server, err := self.SM.Show(id)
	if err != nil {
		utils.Debug(err)
	}
	json, _ = server.JSON()
	return
}

//Del("/api/servers/:id", "delete :id server")
func (self *ServerAPI) Destroy(w http.ResponseWriter, r *http.Request) (json string) {
	return
}

//Put("/api/servers/:id", "update :id server")
func (self *ServerAPI) Update(w http.ResponseWriter, r *http.Request) (json string) {
	return
}



