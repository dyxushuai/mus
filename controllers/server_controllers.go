package controllers

import (
	"github.com/JohnSmithX/mus/manager"
	"net/http"
)

type ServerAPI struct {

}
//Get("/api/servers", "get all")
func (self *ServerAPI) Index(w http.ResponseWriter, r *http.Request) {

}
//Post("/api/servers", "create new")
func (self *ServerAPI) Create(w http.ResponseWriter, r *http.Request) {}

//Get("/api/servers/:id", "get :id server")
func (self *ServerAPI) Show(w http.ResponseWriter, r *http.Request) {}

//Del("/api/servers/:id", "delete :id server")
func (self *ServerAPI) Destroy(w http.ResponseWriter, r *http.Request) {}

//Put("/api/servers/:id", "update :id server")
func (self *ServerAPI) Update(w http.ResponseWriter, r *http.Request) {}



