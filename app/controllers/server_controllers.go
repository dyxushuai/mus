package controllers

import (
	"net/http"
	"github.com/JohnSmithX/mus/app/utils"
	j "encoding/json"
	"github.com/JohnSmithX/mus/app/models"
)

type ServerAPI struct {
}




//Get("/api/servers", "get all")
func (self *ServerAPI) Index(w http.ResponseWriter, r *http.Request) (json string) {
	servers, _ := SM.All()

	data, _ := j.Marshal(servers)
	json = string(data)

	return

}



//Post("/api/servers", "create new")
func (self *ServerAPI) Create(w http.ResponseWriter, r *http.Request) (json string) {
	opt := &models.Server{}

	decoder := j.NewDecoder(r.Body)


	err := decoder.Decode(opt)

	server, err := models.New(opt.Port, opt.Method, opt.Password, opt.Limit, opt.Timeout)

	err = SM.Create(server)
	_ = err
	return
}

//Get("/api/servers/:id", "get :id server")
func (self *ServerAPI) Show(w http.ResponseWriter, r *http.Request) (json string) {
	params := r.URL.Query()
	id := params.Get(":id")


	server, err := SM.Show(id)
	if err != nil {
		utils.Debug(err)
	}
	str, _ := server.JSON()

	json = string(str)
	return
}

//Del("/api/servers/:id", "delete :id server")
func (self *ServerAPI) Destroy(w http.ResponseWriter, r *http.Request) (json string) {
	params := r.URL.Query()
	id := params.Get(":id")

	server, err := SM.Delete(id)
	if err != nil {
		utils.Debug(err)
	}
	err = server.Delete()
	if err != nil {
		utils.Debug(err)
	}
	return
}

//Put("/api/servers/:id", "update :id server")
func (self *ServerAPI) Update(w http.ResponseWriter, r *http.Request) (json string) {


	params := r.URL.Query()
	id := params.Get(":id")
	server, err := SM.Show(id)
	if err != nil {
		return

	}
	opt := &models.Server{}

	decoder := j.NewDecoder(r.Body)
	err = decoder.Decode(opt)

	server, err = models.New(opt.Port, opt.Method, opt.Password, opt.Limit, opt.Timeout)
	err = server.Update()
	err = SM.Create(server)
	_ = err
	return
}



