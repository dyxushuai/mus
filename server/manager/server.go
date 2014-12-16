//this is a shadowsocks server
package manager

import (
	"net"
	"log"
	"sync"
)

const (
	NULL int = iota
	RUN
	STOP
)

type server struct {
	sync.Mutex
	net.Listener
	port string
	state int
	local map[string]*local //1 to 1 : user.username -> local, username must be uniqueness
}

func newServer(port string) (ss *server,err error) {
	if port == "" {
		err = newError("Cannot create a server without port")
		return
	}

	ln, err := net.Listen("tcp", ":" + port)

	if err != nil {
		log.Printf("Create new server for ss at port: %s failed, err: %v\n", port, err)
		err =  err
		return
	}
	ss = &server{sync.Mutex{}, ln, port, STOP, nil}
	return
}


func (self *server) addLocal(user *User) (local *local, err error) {

	local, err = newLocal(user)
	if err != nil {
		return
	}
	self.Lock()
	self.local[user.username] = local
	self.Unlock()
	return
}

func (self *server) close() (err error) {
	err = self.Close()
	return
}

func (self *server) isRun() bool {
	return self.state == RUN
}

func (self *server) run() (err error) {
	if self.isRun() {
		err = newError("Server at port: " + self.port + "is running")
		return
	}
	self.state = RUN

	defer func() {
		err = newError("Server at port: " + self.port + "unexpected stop")
		self.state = STOP
	}()

	for {
		//listen the connect from client
		var conn net.Conn
		conn, err = self.Accept()
		if err != nil {
			log.Printf("accept error: %v\n", err)
			continue
		}

		var user *User
		user, err = getUserFormConn(conn)
		if err != nil {
			log.Printf("Error get passeoed: %s %v\n", self.port, err)
			conn.Close()
			continue
		}

		if user.currentFlow < user.limit {
			log.Printf("Error user runover: %s\n", user.limit)
			conn.Close()
			continue
		}
		//create new client and return
		var local *local
		local, err = self.addLocal(user)
		if err != nil {
			log.Printf("Add client failed: %s %v\n", self.port, err)
			conn.Close()
			continue
		}

		go local.run()
	}
	return
}
