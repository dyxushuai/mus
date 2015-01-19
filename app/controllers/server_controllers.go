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


//func (self *Manager) CreateServerFromBody(body io.Reader) (server models.IServer, err error) {
//	decoder := json.NewDecoder(body)
//	err = decoder.Decode(server)
//	if err != nil {
//		err = utils.NewError(err.Error())
//		return
//	}
//	err = server.InitServer()
//	if err != nil {
//		return
//	}
//	err = self.AddServerToManager(server)
//	if err != nil {
//		return
//	}
//	err = self.AddServerToRedis(server)
//	return
//}
//Post("/api/servers", "create new")
func (self *ServerAPI) Create(w http.ResponseWriter, r *http.Request) (json string) {
	server := &models.Server{}

	decoder := j.NewDecoder(r.Body)


	err := decoder.Decode(server)

	server.Initialize(Store)

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
	return
}



