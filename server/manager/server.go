//this is a shadowsocks server
package manager

import (
	"net"
	"log"
)

type server struct {
	listener net.Listener
	client map[string]*client //1 to 1 : user.username -> client, username must uniqueness
}

func newServer(port string) (*server, error) {
	if port == "" {
		return nil, newError("Cannot create a server without port")
	}
	if ln, err := net.Listen("tcp", ":" + port); err != nil {
		log.Printf("Create new server for ss at port: %s failed, err: %v\n", port, err)
		return nil, err
	} else {
		return &server{listener: ln}, nil
	}
}

//check this user is connect
func (self *server) hasClient(user *User) bool {
	_, ok := self.client[user.username]
	return ok
}

func (self *server) addClient(user *User) (*client, error) {
	if self.hasClient(user) {
		//TODO: add contact method
		return nil, newError("Account is using, please contact the administrator")
	}
	if client, err := newClient(user); err != nil {
		return nil, err
	} else {
		self.client[user.username] = client
		return client, nil
	}
}

func (ss *server) close() error {
	return ss.listener.Close()
}
