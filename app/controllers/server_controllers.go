package controllers

import (
	"net/http"
	j "encoding/json"
	"github.com/JohnSmithX/mus/app/models"
	"github.com/dropbox/godropbox/errors"
)

type ServerAPI struct {
}




//Get("/api/servers", "get all")
func (self *ServerAPI) Index(w http.ResponseWriter, r *http.Request) (json string, err error) {
	servers, err := SM.All()
	if err != nil {
		return
	}
	data, err := j.Marshal(servers)
	if err != nil {
		return
	}
	json = string(data)

	return

}



//Post("/api/servers", "create new")
func (self *ServerAPI) Create(w http.ResponseWriter, r *http.Request) (json string, err error) {
	
	opt := &models.Server{}

	decoder := j.NewDecoder(r.Body)


	err = decoder.Decode(opt)
	if err != nil {
		err = errors.Newf("get params from request failed: %v", err)
		return
	}

	server, err := models.New(opt.Port, opt.Method, opt.Password, opt.Limit, opt.Timeout)
	if err != nil {
		return
	}
	
	err = SM.Create(server)
	if err != nil {
		return
	}
	return
}

//Get("/api/servers/:id", "get :id server")
func (self *ServerAPI) Show(w http.ResponseWriter, r *http.Request) (json string, err error) {
	params := r.URL.Query()
	id := params.Get(":id")


	server, err := SM.Show(id)
	if err != nil {
		return
	}
	data, err := server.JSON()
	if err != nil {
		return
	}
	json = string(data)
	return
}

//Del("/api/servers/:id", "delete :id server")
func (self *ServerAPI) Destroy(w http.ResponseWriter, r *http.Request) (json string, err error) {
	params := r.URL.Query()
	id := params.Get(":id")

	server, err := SM.Delete(id)
	if err != nil {
		return
	}
	err = server.Delete()
	if err != nil {
		return
	}
	return
}

//Put("/api/servers/:id", "update :id server")
func (self *ServerAPI) Update(w http.ResponseWriter, r *http.Request) (json string, err error) {


	params := r.URL.Query()
	id := params.Get(":id")
	server, err := SM.Show(id)
	if err != nil {
		return
	}
	opt := &models.Server{}

	decoder := j.NewDecoder(r.Body)
	err = decoder.Decode(opt)
	if err != nil {
		return
	}
	server, err = models.New(opt.Port, opt.Method, opt.Password, opt.Limit, opt.Timeout)
	if err != nil {
		return
	}
	err = server.Update()
	if err != nil {
		return
	}
	err = SM.Create(server)
	if err != nil {
		return
	}
	return
}



