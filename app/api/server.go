package api

import (
	"github.com/JohnSmithX/mus/db"
	"errors"
)

const (
	serverPrefix = "server:"
	flowPrefix = "flow:"
)

type Server struct {
	Id 				string			`json:"id"`
	Port 			string			`json:"port"`
	Method        	string       	`json:"method"`
	Password      	string       	`json:"password"`
	Limit         	int64        	`json:"limit"`
	Timeout       	int64        	`json:"timeout"`
	Current       	int64       	`json:"current"`
	Create			db.JsonTime		`json:"create_at"`
	Update			db.JsonTime		`json:"update_at"`
}



var store = db.NewStorage()

//operate servers from redis
func  GetServerFromRedis(port string) (server *Server, err error) {
	server, err =  self.store.GetServer(serverPrefix + port)
	if err != nil {
		return
	}
	size, err := self.store.GetSize(flowPrefix + port)
	if err == nil {
		server.current = 0
	} else {
		server.current = size
	}
	err = server.initServer(self)
	return
}

func (self *Server) getServersFromRedis(ports ...string) (servers []*Server, err error) {
	if len(ports) == 0 {
		err = errors.New("Need port but port is nil")
		return
	}
	for _, port := range ports {
		if server, er := self.getServerFromRedis(string(port)); er == nil {
			servers = append(servers, server)
			Debug(er)
		}
	}
	return
}

func (self *Server) getAllServersFromRedis() (servers []*Server, err error) {
	servers, err =  self.store.GetServers(serverPrefix + "**")
	if err != nil {
		return
	}
	for _, server := range servers {
		err = server.initServer(self)
	}
	return
}

func (self *Server) addServerToRedis(server *Server) (err error) {
	err = self.store.SetServer(serverPrefix + server.Port, server)
	return
}

func (self *Server) addServersToRedis(servers []*Server) (err error) {
	for _, server := range servers {
		err = self.addServerToRedis(server)
	}
	return
}

func (self *Server) delServerFromRedis(port string) (err error) {
	err =  self.store.DelServer(serverPrefix + port)
	return
}

func (self *Server) delServersFromRedis(ports ...string) (err error) {
	if len(ports) == 0 {
		err = errors.New("Need port but port is nil")
		return
	}
	for _, port := range ports {
		er := self.delServerFromRedis(port)
		Debug(er)
	}
	return
}

func (self *Server) delAllServersFromRedis() (err error) {
	keys, err := self.store.Keys(serverPrefix + "**")
	if err != nil {
		return
	}
	for _, key := range keys {
		er := self.store.DelServer(key)
		Debug(er)
	}
	return
}
