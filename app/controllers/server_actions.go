package controllers

import (
	"net/http"
)

type ServerActionsAPI struct {

}

//Post("/api/servers/:id/start", "start :id server")
func (self *ServerActionsAPI) Start(w http.ResponseWriter, r *http.Request) {}

//Post("/api/servers/:id/stop", "stop :id server")
func (self *ServerActionsAPI) Stop(w http.ResponseWriter, r *http.Request) {}

//server.Post("/api/servers/:id/restart", "restart :id server")
func (self *ServerActionsAPI) Restart(w http.ResponseWriter, r *http.Request) {}

//server.Get("/api/servers/:id/logs", "get :id server logs")
func (self *ServerActionsAPI) Log(w http.ResponseWriter, r *http.Request) {}

//server.Get("/api/servers/:id/flow", "get :id server flow")
func (self *ServerActionsAPI) Flow(w http.ResponseWriter, r *http.Request) {}
