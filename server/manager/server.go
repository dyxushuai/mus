//this is a shadowsocks server
package manager

import (
	"net"
	"sync"
	"fmt"

)

//server state
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
	comChan ComChan
	local map[string]*local //1 to 1 : user.username -> local, username must be uniqueness
	errFormat string
}

const (
	serverErrFormat string = "Shawdowsocks at port %s : %%s %%v"
)


func newServer(port string) (ss *server,err error) {

	if port == "" {
		err = newError("Cannot create a server without port")
		return
	}

	errFormat := fmt.Sprintf(serverErrFormat, port)

	ln, er := net.Listen("tcp", ":" + port)
	if er != nil {
		err = newError(errFormat, "create listner error:", er)
		return
	}

	ss = &server{sync.Mutex{}, ln, port, STOP, make(ComChan, theNumOfCom), nil, errFormat}
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
	//first stop the loop
	//second close the chan
	//third close the listener
	self.comChan <- CLOSE
	close(self.comChan)
	if er := self.Close(); er != nil {
		err = newError(self.errFormat, "close with error:", er)
	}
	return
}

func (self *server) isRUN() bool {
	return self.state == RUN
}

func (self *server) isSTOP() bool {
	return self.state == STOP
}

func (self *server) isNULL() bool {
	return self.state == NULL
}

func (self *server) listen() (err error) {
//	if !self.isNULL() {
//		err = newError(self.errFormat, "has exist", "")
//		return
//	}
	self.state = RUN

	defer func() {
		err = newError(self.errFormat, "has stoped", "")
		self.state = STOP
	}()
loop:
	for {
		if err != nil {
			bd.addError(err)
		}
		select{
		case com := <- self.comChan:
			switch com {
			case WAIT:
				if !self.isSTOP() {
					self.state = STOP
				}
				continue
			case START:
				if !self.isRUN() {
					self.state = RUN
				}
			case CLOSE:
				break loop
			}
		default:
		}

		conn, er := self.Accept()
		if er != nil {
			err = newError(self.errFormat, "listener accpet error:", er)
			continue
		}

		go self.handleConnect(conn)
	}
	return
}

func(self *server) handleConnect(conn net.Conn) {

	defer func() {
		if err := recover(); err!=nil {
			if v, ok := err.(error); ok {
				fmt.Println(v)
				conn.Write([]byte(v.Error()))
				conn.Close()
			}

		}
	}()

	user, er := getUserFormConn(conn)
	if er != nil {
		panic(newError(self.errFormat, "get user error:", er))
	}
	//create new client and return
	local, er := self.addLocal(user)
	if er != nil {
		panic(newError(self.errFormat, "create connect error:", er))
	}
	local.run()

}
