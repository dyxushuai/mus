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
}

var Log Verbose


func NewManager(verbose bool) (manager *Manager) {

	Log = Verbose(verbose)
	manager = &Manager{}
	manager.servers = make(map[string]*Server)
	return
}


//wrap lock method
func (self *Manager) doWithLock(fn func()) {
	self.mu.Lock()
	defer self.mu.Unlock()
	fn()
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
			Debug(er)
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
	if self.hasServer(server.Port) {
		err = newError("Add proxy server to manager failed: proxy server has existed on port: %s", server.Port)
		return
	}
	self.doWithLock(func () {
		self.servers[server.Port] = server
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
	err = server.destroy()
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
	servers, er = self.getServersFromManager(ports...)
	Debug(er)
	for _, port := range ports {
		_, er = self.delServerFromManager(port)
		Debug(er)
	}
	return
}

func (self *Manager) delAllServersFromManager() (servers []*Server, err error) {


	for port, _ := range self.servers {
		server, er := self.delServerFromManager(port)
		servers = append(servers, server)
		Debug(er)
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
