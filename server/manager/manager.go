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
//var maxOfSS chan map[string]Command




type Manager struct {
	errChan chan error
	sync.Mutex
	ssServers map[string]*server //port -> ss server
}

func CreateManager() (manager *Manager) {
	//default create a chan to receive error message
	errChan := make(chan error, 10)
	manager = &Manager{}
	manager.errChan = errChan
	return
}

//private method for Manager instance
func (self *Manager) hasServer(port string) bool {
	_, ok := self.ssServers[port]
	return ok
}


func (self *Manager) RunAllOfServer() (err []error) {
	for port, _ := range self.ssServers {
		er := self.runServer(port)
		if er != nil {
			err = append(err, er)
		}
	}
	return
}

func (self *Manager) AddServerAndRun(port string) (err error) {
	err = self.AddServer(port)
	if err != nil {
		return
	}
	err = self.runServer(port)
	return
}


//create a new listener with a given port
//each listener with a new goroutine
func (self *Manager) AddServer(port string) (err error) {
	var ss *server
	ss, err = newServer(port)

	if err != nil {
		log.Printf("Create new server for ss at port: %s failed, err: %v\n", port, err)
		return
	}

	if !self.hasServer(port) {
		self.Lock()
		self.ssServers[port] = ss
		self.Unlock()
	} else {
		return newError("Application has listened this port " + port)
	}
	return nil
}

func (self *Manager) runServer(port string) (err error) {
	if !self.hasServer(port) {
		err = newError("No server listen on port: " + port)
		return
	}
	go func() {
		err = self.ssServers[port].run()
		self.errChan <- err
	}()
	return
}

//drop a existed listener
func (self *Manager) DropServer(port string) (err error) {
	if !self.hasServer(port) {
		err = newError(port + " server, drop failed")
		return
	}
	//cannot delete the server which close failed
	err = self.ssServers[port].close()
	if err != nil {
		return
	}
	self.Lock()
	delete(self.ssServers, port)
	self.Unlock()
	return
}




