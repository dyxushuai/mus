//a wrapped manager for ss servers and client

package manager

import (
	"sync"
	"fmt"
)

type Manager struct {
	sync.Mutex
	ssServers map[string]*server //port -> ss server
}


//some command
type Command int
const (
	NIL Command = iota
	WAIT
	START
	CLOSE
)
type ComChan chan Command
var theNumOfCom int = 10
//broadcast
var bd = NewBroadcast()



func CreateManager() (manager *Manager) {
	//default create a chan to receive error message
	manager = &Manager{}
	manager.ssServers = make(map[string]*server)
	return
}



//private method for Manager instance
func (self *Manager) hasServer(port string) bool {
	_, ok := self.ssServers[port]
	return ok
}

func (self *Manager) getServer(port string) (ss *server, err error) {
	if !self.hasServer(port) {
		err = newError("No server listen on the port: %s", port)
		return
	}
	ss = self.ssServers[port]
	return
}



func (self *Manager) runServer(port string) (err error) {
	var ss *server
	ss, err = self.getServer(port)
	if err != nil {
		return
	}
	go func() {
		err = ss.listen()
		if err != nil {
			bd.addError(err)
		}
	}()
	return
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
	if _, er := self.getServer(port); er == nil {
		err = newError("Shadowsocks at port: %s has existed", port)
		return
	}
	ss, er := newServer(port)

	if er != nil {
		err = er
		return
	}
	self.Lock()
	self.ssServers[port] = ss
	self.Unlock()
	return
}


func (self *Manager) StartServer(port string) (err error) {
	ss, er := self.getServer(port)
	if er != nil {
		return
	}
	ss.comChan <- START
	return
}

func (self *Manager) StopServer(port string) (err error) {
	var ss *server
	ss, err = self.getServer(port)
	if err != nil {
		return
	}
	ss.comChan <- WAIT
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
		fmt.Println(<- bd.errChan)
	}
}



