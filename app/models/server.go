package models

import (
	"github.com/JohnSmithX/mus/app/utils"
	"github.com/JohnSmithX/mus/app/db"
	ss "github.com/JohnSmithX/mus/app/shadowsocks"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"time"
	"strings"
	"fmt"

)

//for redis key string
const (
	serverPrefix = "mus:server:"
	flowPrefix = "mus:flow:"
)

type Server struct {

	store 				db.IStorage
	Id 					string			`json:"id"`
	CreateTime			utils.Time		`json:"create_at"`
	UpdateTime			utils.Time		`json:"update_at"`
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



	server.InitServer(recorder)

	err = server.save()
	return
}

func addPrefix(key, prefix string) string {
	if strings.HasPrefix(key, prefix) {
		return key
	}
	return prefix + key
}
//json ID
func (self *Server) Initialize(recorder db.IStorage) (err error) {
	err = self.Server.InitServer(recorder)
	self.Id = uuid.NewV4().String()
	self.store = recorder
	self.upTime()
	self.crTime()

	return
}

//update time at Now
func (self *Server) upTime() {
	self.UpdateTime = utils.Time(time.Now())
}

//create time at Now
func (self *Server) crTime() {
	self.CreateTime = utils.Time(time.Now())
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

func (self *Server) JSON() (result []byte, err error) {
	data, err := json.Marshal(self)
	if err != nil {
		return
	}
	result = data
	return
}


//operate servers from redis
func GetServerFromRedis(store db.IStorage, port string) (server *Server, err error) {
	data, err :=  store.GetServer(addPrefix(port, serverPrefix))

	if err != nil {

		return
	}

	server = &Server{}
	err = json.Unmarshal(data, server)

	if err != nil {

		return
	}

	size, err := store.GetSize(addPrefix(port, flowPrefix))
	if err != nil {
		store.IncrSize(addPrefix(port, flowPrefix), 0)
		server.Current = 0
	} else {

		server.Current = size
	}


	server.Initialize(store)

	fmt.Println(err)
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
		} else {
			err = er
			return
		}
	}
	return
}

func GetAllServersFromRedis(store db.IStorage) (servers []*Server, err error) {
	keys, err := store.Keys(serverPrefix + "**")


	if err != nil {
		return
	}
	for _, key := range keys {

		if server, er := GetServerFromRedis(store, key); er == nil {
			servers = append(servers, server)
		} else {
			err = er

			return
		}
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
		if err != nil {
			return
		}
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
		err = DelServerFromRedis(store, string(port))
		if err != nil {
			return
		}
	}
	return
}

func DelAllServersFromRedis(store db.IStorage, ) (err error) {
	keys, err := store.Keys(serverPrefix + "**")
	if err != nil {
		return
	}
	for _, key := range keys {
		err = store.DelServer(key)
		if err != nil {
			return
		}
	}
	return
}
