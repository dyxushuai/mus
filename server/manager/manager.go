//a wrapped manager for ss servers and client

package manager

import (
	"log"
	"sync"


)



var theNumOfErr int = 10

type Manager struct {
	*sync.Mutex
	errChan chan error//error message channel
	ssServers map[string]*server //port -> ss server
}

func CreateManager() (manager *Manager) {
	//default create a chan to receive error message

	errChan := make(chan error, theNumOfErr)
	manager = &Manager{}
	manager.errChan = errChan
	go manager.LogError()
	return
}

//this function is for log all of error message
func (self *Manager) LogError() {
	for {
		select {
		case e := <- self.errChan:
			log.Println(e)
		}
	}
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
		return
	}

	if !self.hasServer(port) {
		self.Lock()
		self.ssServers[port] = ss
		self.Unlock()
	} else {
		err = newError("Application has listened this port " + port)
	}
	return
}

func (self *Manager) runServer(port string) (err error) {
	if !self.hasServer(port) {
		err = newError("No server listen on the port: " + port)
		return
	}
	go func() {
		err = self.ssServers[port].run()
		self.errChan <- err
	}()
	return
}

func (self *Manager) StartServer(port string) (err error) {
	if !self.hasServer(port) {
		err = newError("No server Listen on the port: " + port)
		return
	}
}

func (self *Manager) StopServer(port string) (err error) {
	if !self.hasServer(port) {
		err = newError("No server Listen on the port: " + port)
		return
	}
	self.ssServers[port].comChan <- WAIT
	return
}

//drop a existed listener
func (self *Manager) DropServer(port string) (err error) {
	if !self.hasServer(port) {
		err = newError("No server Listen on the port: " + port)
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




