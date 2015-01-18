package models

import (
	"github.com/JohnSmithX/mus/app/utils"
	"github.com/JohnSmithX/mus/app/db"
	ss "github.com/JohnSmithX/mus/app/shadowsocks"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"time"

)

//for redis key string
const (
	serverPrefix = "server:"
	flowPrefix = "flow:"
)

type Server struct {
	ss.Server
	store 			db.IStorage
	Id 				string				`json:"id"`
	Create			utils.JsonTime		`json:"create_at"`
	Update			utils.JsonTime		`json:"update_at"`
}

func New(port, method, password string, limit, timeout int64 ,recorder db.IStorage) (server *Server,err error) {
	if port == "" {
		err = utils.NewError("Cannot create a server without port")
		return
	}

	server = &Server{}
	server.Server, err = ss.NewServer(port, method, password, limit, timeout, recorder)
	if err != nil {
		return
	}

	server.store = recorder

	server.initialize()

	err = server.InitServer()
	
	err = server.save()
	return
}

//json ID
func (self *Server) initialize() {
	self.Id = uuid.NewV4()
	self.upTime()
	self.crTime()
}

//update time at Now
func (self *Server) upTime() {
	self.Update.Time = time.Now()
}

//create time at Now
func (self *Server) crTime() {
	self.Create.Time = time.Now()
}


func (self *Server) save() (err error) {
	err = AddServerToRedis(self.store, self)
	return
}


func (self *Server) Update() (err error) {
	self.upTime()
	err = self.save()
	return
}

func (self *Server) Delete() (err error) {
	err = DelServerFromRedis(self.store, self.Port)
	return
}


//operate servers from redis
func GetServerFromRedis(store db.IStorage, port string) (server *Server, err error) {
	data, err :=  store.GetServer(serverPrefix + port)

	if err != nil {
		return
	}
	err = json.Unmarshal(data, server)
	if err != nil {
		return
	}

	size, err := store.GetSize(flowPrefix + port)
	if err == nil {
		server.Current = 0
	} else {
		server.Current = size
	}
	server.recorder = store

	err = server.InitServer()
	return
}

func GetServersFromRedis(store db.IStorage, ports ...string) (servers []*Server, err error) {
	if len(ports) == 0 {
		err = utils.NewError("Need port but port is nil")
		return
	}
	for _, port := range ports {
		if server, er := GetServerFromRedis(store, string(port)); er == nil {
			servers = append(servers, server)
			utils.Debug(er)
		}
	}
	return
}

func GetAllServersFromRedis(store db.IStorage) (servers []*Server, err error) {
	servers, err =  store.GetServers(serverPrefix + "**")
	if err != nil {
		return
	}
	for _, server := range servers {
		err = server.InitServer()
	}
	return
}

func AddServerToRedis(store db.IStorage, server *Server) (err error) {
	data, err := json.Marshal(server)
	err = store.SetServer(serverPrefix + server.Port, data)
	return
}

func AddServersToRedis(store db.IStorage, servers []*Server) (err error) {
	for _, server := range servers {
		err = AddServerToRedis(store, server)
	}
	return
}

func DelServerFromRedis(store db.IStorage, port string) (err error) {
	err =  store.DelServer(serverPrefix + port)
	return
}

func DelServersFromRedis(store db.IStorage, ports ...string) (err error) {
	if len(ports) == 0 {
		err = utils.NewError("Need port but port is nil")
		return
	}
	for _, port := range ports {
		er := DelServerFromRedis(store, string(port))
		utils.Debug(er)
	}
	return
}

func DelAllServersFromRedis(store db.IStorage, ) (err error) {
	keys, err := store.Keys(serverPrefix + "**")
	if err != nil {
		return
	}
	for _, key := range keys {
		er := store.DelServer(key)
		utils.Debug(er)
	}
	return
}
