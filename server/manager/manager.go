//a wrapped manager for ss servers and client

package manager

import (
	"log"
	"sync"


)


//some command
type Command int

//const (
//	NULL Command = iota
//	WAIT Command
//	OPEN Command
//	CLOSE Command
//)
//default commands channel with 10 maximum
var maxOfSS chan map[string]Command

type Manager struct {
	sync.Mutex
	ssServers map[string]*server //port -> ss server
}

func CreateManager() *Manager {
	return &Manager{}
}

//private method for Manager instance
func (self *Manager) hasServer(port string) bool {
	_, ok := self.ssServers[port]
	return ok
}

//
//func (self *Manager) addClient(port string, user *User) (*client, error) {
//	if self.hasServer(port) {
//		return self.ssServers[port].addClient(user)
//	}
//	return nil, newError("no shadowsocks server listen on " + port)
//}

//create a new listener with a given port
//each listener with a new goroutine
func (self *Manager) AddServerAndRun(port string) error {
	if server, err := newServer(port); err != nil {
		log.Printf("Create new server for ss at port: %s failed, err: %v\n", port, err)
		return err
	} else {
		if !self.hasServer(port) {
			self.Lock()
			self.ssServers[port] = server
			self.Unlock()
		} else {
			return newError("Application has listened this port " + port)
		}
	}
	return nil
}

//drop a existed listener
func (self *Manager) DropServer(port string) error {
	if !self.hasServer(port) {
		return newError(port + " server inexistence, drop failed")
	}
	self.Lock()
	err := self.ssServers[port].close()
	delete(self.ssServers, port)
	self.Unlock()
	return err
}




