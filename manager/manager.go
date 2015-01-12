//a wrapped manager for ss servers and client

package manager

import (
	"sync"
	"github.com/JohnSmithX/mus/config"
//	"fmt"
)

type Manager struct {
	mu sync.Mutex
	servers map[string]*Server //port -> ss server
	store *Storage
}


var log Verbose

//create redis connect pool
var storage *Storage


func init() {
	storage = NewStorage(config.REDIS_SERVER, config.REDIS_PASSWORD)
	if err := storage.Test(); err != nil {
		panic("connect to redis server error: " + err.Error())
	}
}



func New(verbose bool) (manager *Manager) {

	log = Verbose(verbose)
	manager = &Manager{}
	manager.servers = make(map[string]*Server)
	manager.store = storage
	return
}



//wrap lock method
func (self *Manager) withLockDo(fn func()) {
	self.mu.Lock()
	defer self.mu.Unlock()
	fn()
}

//private method for Manager instance
func (self *Manager) hasServer(port string) bool {
	_, ok := self.servers[port]
	return ok
}

func (self *Manager) getServer(port string) (ss *Server, err error) {
	if !self.hasServer(port) {
		err = newError("Thers is no proxy server listened on the port: %s", port)
		return
	}
	ss = self.servers[port]
	return
}



func (self *Manager) AddServerAndStart(port, method, password string, limit, timeout int64) (err error) {
	defer func() {
		if err != nil {
			log.Debug(err.Error())
		}
	}()

	id, err := self.AddServer(port, method, password, limit, timeout)
	if err != nil {
		return
	}
	err = self.StartServer(id)
	return
}

//create a new listener with a given args
func (self *Manager) AddServer(port, method, password string, limit, timeout int64) (id string, err error) {
	defer func() {
		if err != nil {
			log.Debug(err.Error())
		}
	}()

	if self.hasServer(port) {
		err = newError("Add proxy server failed: proxy server has listened on port %s", port)
		return
	}

	ss, er := newServer(port, method, password, limit, timeout, self.store)

	if er != nil {
		err = er
		return
	}
	id = port

	self.withLockDo(func() {
		self.servers[port] = ss
	})
	return
}

func (self *Manager) StartServer(port string) (err error) {
	defer func() {
		if err != nil {
			log.Debug(err.Error())
		}
	}()

	if !self.hasServer(port) {
		err = newError("Start proxy server failed: no server listen on port %s", port)
		return
	}
	err = self.servers[port].Start()
	return
}
//stop a started server
func (self *Manager) StopServer(port string) (err error) {
	defer func() {
		if err != nil {
			log.Debug(err.Error())
		}
	}()

	ss, err := self.getServer(port)
	if err != nil {
		return
	}
	ss.Stop()
	return
}

//drop a existed listener
func (self *Manager) DropServer(port string) (err error) {
	defer func() {
		if err != nil {
			log.Debug(err.Error())
		}
	}()


	ss, err := self.getServer(port)
	if err != nil {
		return
	}

	err = ss.Close()
	//cannot delete the server which close failed
	if err != nil {
		return
	}

	self.withLockDo(func() {
		delete(self.servers, port)
	})
	return
}



