//a wrapped manager for ss servers and client
//port is the primary key
package manager

import (
	"sync"
	"encoding/json"
	"io"
)

type Manager struct {
	mu sync.Mutex
	verbose bool
	servers map[string]*Server //port -> ss server
	store *Storage
}

var log Verbose


func New(host, password string, verbose bool) (manager *Manager, err error) {

	log = Verbose(verbose)
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

//private method for Manager instance
func (self *Manager) hasServer(port string) bool {
	_, ok := self.servers[port]
	return ok
}


//operate servers from manager
func (self *Manager) getServerFromManager(port string) (server *Server, err error) {
	if !self.hasServer(port) {
		err = newError("There is no proxy server listened on the port: %s", port)
		return
	}
	server = self.servers[port]
	return
}

func (self *Manager) getServersFromManager(ports ...string) (servers []*Server, err error) {
	if len(ports) == 0 {
		err = newError("Need port but port is nil")
		return
	}
	for _, port := range ports {
		if server, _ := self.getServerFromManager(port); server != nil {
			servers = append(servers, server)
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

//operate servers from redis
func (self *Manager) getServerFromRedis(port string) (server *Server, err error) {
	server, err =  self.store.GetServer(serverPrefix + port)
	if err != nil {
		return
	}
	err = server.initServer()
	return
}

func (self *Manager) getServersFromRedis(ports ...string) (servers []*Server, err error) {
	if len(ports) == 0 {
		err = newError("Need port but port is nil")
		return
	}
	for _, port := range ports {
		if server, err := self.getServerFromRedis(port); err == nil {
			servers = append(servers, server)
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
		err = server.initServer()
	}
	self.addServersToManager(servers)
	return
}

func (self *Manager) addServerToRedis(server *Server) (err error) {
	err = self.store.SetServer("server:" + server.Port, server)
	return
}

func (self *Manager) addServersToRedis(servers []*Server) (err error) {
	for _, server := range servers {
		err = self.addServerToRedis(server)
	}
	return
}

//API
//request with json content
func (self *Manager) CreateServerFromBody(body io.Reader) (server *Server, err error) {
	decoder := json.NewDecoder(body)
	err = decoder.Decode(server)
	if err != nil {
		err = newError(err.Error())
		return
	}
	err = server.initServer()
	return
}


func (self *Manager) StartServer(port string) (err error) {
	defer func() {
		if err != nil {
			log.Debug(err.Error())
		}
	}()

	if server, err := self.getServerFromManager(port); err == nil {
		err = server.Start()
	}
	return
}
//stop a started server
func (self *Manager) StopServer(port string) (err error) {
	defer func() {
		if err != nil {
			log.Debug(err.Error())
		}
	}()

	if server, err := self.getServerFromManager(port); err == nil {
		err = server.Stop()
	}
	return
}

//drop a existed listener
func (self *Manager) DestroyServer(port string) (err error) {
	defer func() {
		if err != nil {
			log.Debug(err.Error())
		}
	}()


	server, err := self.getServerFromManager(port)
	if err != nil {
		return
	}

	err = server.Destroy()
	//cannot delete the server which close failed
	if err != nil {
		return
	}

	self.doWithLock(func() {
		delete(self.servers, port)
	})
	return
}



