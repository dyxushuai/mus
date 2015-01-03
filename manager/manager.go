//a wrapped manager for ss servers and client

package manager

import (
	"sync"
	"github.com/JohnSmithX/mus/config"
)

type Manager struct {
	mu sync.Mutex
	ssServers map[string]*server //port -> ss server
}

type config interface {
	Config() (string, string, string, int64)
}

type servers []server

//broadcast
var bd = NewBroadcast()

var redis *Storage


func New() (manager *Manager) {

	//create redis connect pool
	redis = config.NewStorage()

	manager = &Manager{}
	manager.ssServers = make(map[string]*server)

	return
}

//wrap lock method
func (self *Manager)withLockDo(fn interface {}) {
	
}

//private method for Manager instance
func (self *Manager) hasServer(port string) bool {
	_, ok := self.ssServers[port]
	return ok
}

func (self *Manager) getServer(port string) (ss *server, err error) {
	if !self.hasServer(port) {
		err = newError("Thers is no proxy server listened on the port: %s", port)
		bd.addError(err)
		return
	}
	ss = self.ssServers[port]
	return
}


func (self *Manager) StartServer(port string) (err error) {

	if !self.hasServer(port) {
		err = newError("Start proxy server failed: no server listen on port %s", port)
		bd.addError(err)
		return
	}
	err = self.ssServers[port].start()
	if err != nil {
		bd.addError(err)
	}
	return
}

//current proxy server list
func (self *Manager) ServerList() servers {
	var list servers
	for _, v := range self.ssServers {
		list = append(list, v)
	}
	return list
}

////run all of proxy which hasn't started
//func (self *Manager) RunAllOfServer() (err []error) {
//	for port, _ := range self.ssServers {
//		er := self.startServer(port)
//		if er != nil {
//			err = append(err, er)
//		}
//	}
//	return
//}

func (self *Manager) AddServerAndRun(conf config) (err error) {

	id, err := self.AddServer(conf)
	if err != nil {
		return
	}
	err = self.StartServer(id)
	return
}


//create a new listener with a given port
//each listener with a new goroutine
func (self *Manager) AddServer(conf config) (id string, err error) {
	port, method, password, timeout := conf.Config()


	if self.hasServer(port) {
		err = newError("Add proxy server failed: proxy server has listened on port %s", port)
		bd.addError(err)
		return
	}
	ss, er := newServer(port, method, password, timeout)

	if er != nil {
		err = er
		return
	}
	id = port
	self.Lock()
	self.ssServers[port] = ss
	self.Unlock()
	return
}



//stop a started server
func (self *Manager) StopServer(port string) (err error) {
	var ss *server
	ss, err = self.getServer(port)
	if err != nil {
		return
	}
	ss.stop()
	return
}

//drop a existed listener
func (self *Manager) DropServer(port string) (err error) {
	var ss *server
	ss, err = self.getServer(port)
	if err != nil {
		return
	}
	//cannot delete the server which close failed
	err = ss.close()
	if err != nil {
		return
	}
	self.Lock()
	delete(self.ssServers, port)
	self.Unlock()
	return
}


func (self *Manager) DEBUG() {
	for {
		select {
		case err := <- bd.errChan:
			if v, ok := err.(*errorType); ok {
				v.print()
			}
		}
	}
}



func (self *Manager) LOG() {
	for {
		select {
		case msg := <- bd.msgChan:
			if v, ok := msg.(log); ok {
				v.print()
			}
		}
	}
}


