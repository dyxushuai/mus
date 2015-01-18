//a wrapped manager for ss servers and client
//port is the primary key
package manager

import (
	"sync"
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
func (self *Manager) GetServerFromManager(port string) (server models.IServer, err error) {
	err = self.validServer(port)
	if err != nil {
		return
	}
	self.doWithLock(func () {
		server = self.servers[port]
	})
	return
}

func (self *Manager) GetServersFromManager(ports ...string) (servers []models.IServer, err error) {
	if len(ports) == 0 {
		err = utils.NewError("Need port but port is nil")
		return
	}

	for _, port := range ports {
		if server, er := self.GetServerFromManager(string(port)); er == nil {
			servers = append(servers, server)
		} else {
			err = er
			return
		}
	}
	return
}

func (self *Manager) GetAllServersFromManager() (servers []models.IServer, err error) {
	if len(self.servers) == 0 {
		err = utils.NewError("There is no proxy server in manager")
		return
	}
	for _, server := range self.servers {
		servers = append(servers, server)
	}
	return
}

func (self *Manager) AddServerToManager(server models.IServer) (err error) {
	if self.hasServer(server.Key()) {
		err = utils.NewError("Add proxy server to manager failed: proxy server has existed on port: %s", server.Key())
		return
	}
	self.doWithLock(func () {
		self.servers[server.Key()] = server
	})

	return
}

func (self *Manager) AddServersToManager(servers []models.IServer) (err error) {
	for _, server := range servers {
		err = self.AddServerToManager(server)
	}
	return
}

func (self *Manager) DelServerFromManager(port string) (server models.IServer, err error) {
	err = self.validServer(port)
	if err != nil {
		return
	}
	server = self.servers[port]
	err = server.Destroy()

	if err != nil {
		return
	}

	self.doWithLock(func () {
		delete(self.servers, port)
	})
	return
}

func (self *Manager) DelServersFromManager(ports ...string) (servers []models.IServer, err error) {
	if len(ports) == 0 {
		err = utils.NewError("Need port but port is nil")
		return
	}

	var server models.IServer
	for _, port := range ports {
		server, err = self.DelServerFromManager(string(port))
		if err != nil {
			return
		}
		servers = append(servers, server)
	}
	return
}

func (self *Manager) DelAllServersFromManager() (servers []models.IServer, err error) {

	var server models.IServer
	for port, _ := range self.servers {
		server, err = self.DelServerFromManager(port)
		if err != nil {
			return
		}
		servers = append(servers, server)
	}
	return
}


//TODO: API
//request with json content
//POST /api/servers
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


//GET /api/servers
func (self *Manager) All() (servers []models.IServer, err error)  {
	servers, err = self.GetAllServersFromManager()
	return
}

//GET /api/servers/:id select
func (self *Manager) Show(id string) (server models.IServer, err error) {
	server, err = self.GetServerFromManager(id)
	return
}

//DEL /api/servers/:id delete
func (self *Manager) Delete(id string) (server models.IServer, err error) {
	server, err = self.DelServerFromManager(id)
	return
}

//PUT /api/servers/:id update
//func (self *Manager) Update(id string, body io.Reader) (server models.IServer, err error) {
//	server, err = self.Delete(id)
//	if err != nil {
//		return
//	}
//	server, err = self.CreateServerFromBody(body)
//	return
//}
