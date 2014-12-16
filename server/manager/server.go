//this is a shadowsocks server
package manager

import (
	"net"
	"log"
)

type server struct {
	port string
	listener net.Listener
	local map[string]*local //1 to 1 : user.username -> client, username must uniqueness
}

func newServer(port string) (server *server,err error) {
	if port == "" {
		err = newError("Cannot create a server without port")
	}
	if ln, err := net.Listen("tcp", ":" + port); err != nil {
		log.Printf("Create new server for ss at port: %s failed, err: %v\n", port, err)
		err =  err
	} else {
		server.port = port
		server.listener = ln
	}
	return
}

//check this user is connect
func (self *server) hasLocal(user *User) bool {
	_, ok := self.local[user.username]
	return ok
}

func (self *server) addLocal(user *User) (local *local, err error) {
	if self.hasLocal(user) {
		//TODO: add contact method
		err = newError("Account is using, please contact the administrator")
	}
	local, err = newLocal(user)
	if err != nil {
		return
	}
	self.local[user.username] = local
	return
}

func (self *server) close() (err error) {
	err = self.listener.Close()
	return
}

func (self *server) run() (err error) {
	for {
		//listen the connect from client
		var conn net.Conn
		conn, err = self.listener.Accept()
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
