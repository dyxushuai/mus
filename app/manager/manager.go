//a wrapped manager for ss servers and client
//port is the primary key
package manager

import (
	"sync"
	"encoding/json"
	"io"
	"github.com/JohnSmithX/mus/app/utils"
	"github.com/JohnSmithX/mus/app/models"
)



type Manager struct {
	mu sync.Mutex
	servers map[string]models.IServer //port -> ss server
}



func NewManager() (manager *Manager) {

	manager = &Manager{}
	manager.servers = make(map[string]models.IServer)
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
		err = utils.NewError("There is no proxy server listened on the port: %s", port)
	}
	return
}


//operate servers from manager
func (self *Manager) getServerFromManager(port string) (server models.IServer, err error) {
	err = self.validServer(port)
	if err != nil {
		return
	}
	self.doWithLock(func () {
		server = self.servers[port]
	})
	return
}

func (self *Manager) getServersFromManager(ports ...string) (servers []models.IServer, err error) {
	if len(ports) == 0 {
		err = utils.NewError("Need port but port is nil")
		return
	}

	for _, port := range ports {
		if server, er := self.getServerFromManager(string(port)); er == nil {
			servers = append(servers, server)
			utils.Debug(er)
		}
	}
	return
}

func (self *Manager) getAllServersFromManager() (servers []models.IServer, err error) {
	if len(self.servers) == 0 {
		err = utils.NewError("There is no proxy server in manager")
		return
	}
	for _, server := range self.servers {
		servers = append(servers, server)
	}
	return
}

func (self *Manager) addServerToManager(server models.IServer) (err error) {
	if self.hasServer(server.port) {
		err = utils.NewError("Add proxy server to manager failed: proxy server has existed on port: %s", server.port)
		return
	}
	self.doWithLock(func () {
		self.servers[server.port] = server
	})

	return
}

func (self *Manager) addServersToManager(servers []models.IServer) (err error) {
	for _, server := range servers {
		err = self.addServerToManager(server)
	}
	return
}

func (self *Manager) delServerFromManager(port string) (server models.IServer, err error) {
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

func (self *Manager) delServersFromManager(ports ...string) (servers []models.IServer, err error) {
	if len(ports) == 0 {
		err = utils.NewError("Need port but port is nil")
		return
	}
	var er error
	servers, er = self.getServersFromManager(ports...)
	utils.Debug(er)
	for _, port := range ports {
		_, er = self.delServerFromManager(string(port))
		utils.Debug(er)
	}
	return
}

func (self *Manager) delAllServersFromManager() (servers []models.IServer, err error) {


	for port, _ := range self.servers {
		server, er := self.delServerFromManager(port)
		servers = append(servers, server)
		utils.Debug(er)
	}
	return
}


//TODO: API
//request with json content
//POST /api/servers
func (self *Manager) CreateServerFromBody(body io.Reader) (server models.IServer, err error) {
	decoder := json.NewDecoder(body)
	err = decoder.Decode(server)
	if err != nil {
		err = utils.NewError(err.Error())
		return
	}
	err = server.initServer()
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
func (self *Manager) All() (servers []models.IServer, err error)  {
	servers, err = self.getAllServersFromManager()
	return
}

//GET /api/servers/:id select
func (self *Manager) Show(id string) (server models.IServer, err error) {
	server, err = self.getServerFromManager(id)
	return
}

//DEL /api/servers/:id delete
func (self *Manager) Delete(id string) (server models.IServer, err error) {
	server, err = self.delServerFromManager(id)
	return
}

//PUT /api/servers/:id update
func (self *Manager) Update(id string, body io.Reader) (server models.IServer, err error) {
	server, err = self.Delete(id)
	if err != nil {
		return
	}
	server, err = self.CreateServerFromBody(body)
	return
}
