package api

import (
	"github.com/JohnSmithX/mus/app/utils"
	"github.com/JohnSmithX/mus/app/db"
	"github.com/JohnSmithX/mus/app/manager"
	"encoding/json"

)

const (
	serverPrefix = "server:"
	flowPrefix = "flow:"
)

type Server struct {
	manager.Manager
	Id 				string				`json:"id"`
	Create			utils.JsonTime		`json:"create_at"`
	Update			utils.JsonTime		`json:"update_at"`
}


//operate servers from redis
func GetServerFromRedis(store *db.Storage, port string) (server *Server, err error) {
	server, err =  store.GetServer(serverPrefix + port)

	if err != nil {
		return
	}
	size, err := store.GetSize(flowPrefix + port)
	if err == nil {
		server.current = 0
	} else {
		server.current = size
	}
	err = server.initServer()
	return
}

func GetServersFromRedis(store *db.Storage, ports ...string) (servers []*Server, err error) {
	if len(ports) == 0 {
		err = utils.NewError("Need port but port is nil")
		return
	}
	for _, port := range ports {
		if server, er := GetServerFromRedis(store, string(port)); er == nil {
			servers = append(servers, server)
			Debug(er)
		}
	}
	return
}

func GetAllServersFromRedis(store *db.Storage, ) (servers []*Server, err error) {
	servers, err =  store.GetServers(serverPrefix + "**")
	if err != nil {
		return
	}
	for _, server := range servers {
		err = server.initServer()
	}
	return
}

func  AddServerToRedis(store *db.Storage, server *Server) (err error) {
	data, err := json.Marshal(server)
	err = store.SetServer(serverPrefix + server.Port, data)
	return
}

func  addServersToRedis(store *db.Storage, servers []*Server) (err error) {
	for _, server := range servers {
		err = AddServerToRedis(store, server)
	}
	return
}

func  DelServerFromRedis(store *db.Storage, port string) (err error) {
	err =  store.DelServer(serverPrefix + port)
	return
}

func  DelServersFromRedis(store *db.Storage, ports ...string) (err error) {
	if len(ports) == 0 {
		err = utils.NewError("Need port but port is nil")
		return
	}
	for _, port := range ports {
		er := DelServerFromRedis(store, string(port))
		Debug(er)
	}
	return
}

func  delAllServersFromRedis(store *db.Storage, ) (err error) {
	keys, err := store.Keys(serverPrefix + "**")
	if err != nil {
		return
	}
	for _, key := range keys {
		er := store.DelServer(key)
		Debug(er)
	}
	return
}
