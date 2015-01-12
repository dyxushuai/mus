//a wrapped manager for ss servers and client
//port is the primary key
package manager

import (
	"sync"
	"encoding/json"
	"io"
	"fmt"
)

type Manager struct {
	mu sync.Mutex
	verbose bool
	servers map[string]*Server //port -> ss server
	store *Storage
}

var Log Verbose


func New(host, password string, verbose bool) (manager *Manager, err error) {

	Log = Verbose(verbose)
	manager = &Manager{}
	manager.servers = make(map[string]*Server)

	//create redis connect pool
	manager.store = NewStorage(host, password)
	err = manager.initialize()
	return
}

//initialize when call `New`
//first get all servers from redis
//add them to manager
func (self *Manager) initialize() (err error) {

	servers, err := self.getAllServersFromRedis()
	if err != nil {
		return
	}
	err = self.addServersToManager(servers)
	return
}

//wrap lock method
func (self *Manager) doWithLock(fn func()) {
	self.mu.Lock()
	defer self.mu.Unlock()
	fn()
}

func (self *Manager) Debug(err error) {
	if err !=nil {
		Log.Debug(err.Error())
	}
}

//private method for Manager instance
func (self *Manager) hasServer(port string) bool {
	_, ok := self.servers[port]
	return ok
}

func (self *Manager) validServer(port string) (err error) {
	if !self.hasServer(port) {
		err = newError("There is no proxy server listened on the port: %s", port)
	}
	return
}


//api util
//operate servers from redis
func (self *Manager) getServerFromRedis(port string) (server *Server, err error) {
	server, err =  self.store.GetServer(serverPrefix + port)
	if err != nil {
		return
	}
	err = server.initServer(self)
	return
}

func (self *Manager) getServersFromRedis(ports ...string) (servers []*Server, err error) {
	if len(ports) == 0 {
		err = newError("Need port but port is nil")
		return
	}
	for _, port := range ports {
		if server, er := self.getServerFromRedis(port); er == nil {
			servers = append(servers, server)
			self.Debug(er)
		}
	}
	return
}

func (self *Manager) getAllServersFromRedis() (servers []*Server, err error) {
	servers, err =  self.store.GetServers(serverPrefix + "**")
	if err != nil {
		return
	}
	for _, server := range servers {
		err = server.initServer(self)
	}
	return
}

func (self *Manager) addServerToRedis(server *Server) (err error) {
	err = self.store.SetServer(serverPrefix + server.Port, server)
	return
}

func (self *Manager) addServersToRedis(servers []*Server) (err error) {
	for _, server := range servers {
		err = self.addServerToRedis(server)
	}
	return
}

func (self *Manager) delServerFromRedis(port string) (err error) {
	err =  self.store.DelServer(serverPrefix + port)
	return
}

func (self *Manager) delServersFromRedis(ports ...string) (err error) {
	if len(ports) == 0 {
		err = newError("Need port but port is nil")
		return
	}
	for _, port := range ports {
		er := self.delServerFromRedis(port)
		self.Debug(er)
	}
	return
}

func (self *Manager) delAllServersFromRedis() (err error) {
	keys, err := self.store.Keys(serverPrefix + "**")
	if err != nil {
		return
	}
	for _, key := range keys {
		er := self.store.DelServer(key)
		self.Debug(er)
	}
	return
}

//operate servers from manager
func (self *Manager) getServerFromManager(port string) (server *Server, err error) {
	err = self.validServer(port)
	if err != nil {
		return
	}
	self.doWithLock(func () {
		server = self.servers[port]
	})
	return
}

func (self *Manager) getServersFromManager(ports ...string) (servers []*Server, err error) {
	if len(ports) == 0 {
		err = newError("Need port but port is nil")
		return
	}

	for _, port := range ports {
		if server, er := self.getServerFromManager(port); er == nil {
			servers = append(servers, server)
			self.Debug(er)
		}
	}
	return
}

func (self *Manager) getAllServersFromManager() (servers []*Server, err error) {
	if len(self.servers) == 0 {
		err = newError("There is no proxy server in manager")
		return
	}
	for _, server := range self.servers {
		servers = append(servers, server)
	}
	return
}

func (self *Manager) addServerToManager(server *Server) (err error) {
	self.doWithLock(func() {
		if self.hasServer(server.Port) {
			err = newError("Add proxy server to manager failed: proxy server has existed on port: %s", server.Port)
			return
		}
		self.doWithLock(func () {
			self.servers[server.Port] = server
		})
	})
	return
}

func (self *Manager) addServersToManager(servers []*Server) (err error) {
	for _, server := range servers {
		err = self.addServerToManager(server)
	}
	return
}

func (self *Manager) delServerFromManager(port string) (server *Server, err error) {
	err = self.validServer(port)
	if err != nil {
		return
	}
	server = self.servers[port]
	err = server.Destroy()
	self.doWithLock(func () {
		delete(self.servers, port)
	})
	return
}

func (self *Manager) delServersFromManager(ports ...string) (servers []*Server, err error) {
	if len(ports) == 0 {
		err = newError("Need port but port is nil")
		return
	}

	var er error
	for _, port := range ports {
		servers, er = self.getServersFromManager(port...)
		self.Debug(er)
		er = self.delServerFromManager(port)
		self.Debug(er)
	}
	return
}

func (self *Manager) delAllServersFromManager() (servers []*Server, err error) {

	var er error
	for port, _ := range self.servers {
		servers, er = self.delServerFromManager(port)
		self.Debug(er)
	}
	return
}


//TODO: API
//request with json content
//POST /api/servers
func (self *Manager) CreateServerFromBody(body io.Reader) (server *Server, err error) {
	decoder := json.NewDecoder(body)
	err = decoder.Decode(server)
	if err != nil {
		err = newError(err.Error())
		return
	}
	err = server.initServer(self)
	if err != nil {
		return
	}
	err = self.addServerToManager(server)
	if err != nil {
		return
	}
	err = self.addServerToRedis(server)
	return
}

func (self *Manager) CreateServerFromArgs(port, method, password string, limit, timeout int64) (server *Server, err error) {
	server, err = newServer(port, method, password, limit, timeout, self)
	if err != nil {
		return
	}
	fmt.Println("here")
	err = self.addServerToManager(server)

	if err != nil {
		return
	}
	err = self.addServerToRedis(server)
	return
}

//GET /api/servers
func (self *Manager) All() (servers []*Server, err error)  {
	servers, err = self.getAllServersFromManager()
	return
}

//GET /api/servers/:id select
func (self *Manager) Show(id string) (server *Server, err error) {
	server, err = self.getServerFromManager(id)
	return
}

//DEL /api/servers/:id delete
func (self *Manager) Delete(id string) (server *Server, err error) {
	server, err = self.delServerFromManager(id)
	return
}

//PUT /api/servers/:id update
func (self *Manager) Update(id string, body io.Reader) (server *Server, err error) {
	server, err = self.Delete(id)
	if err != nil {
		return
	}
	server, err = self.CreateServerFromBody(body)
	return
}
