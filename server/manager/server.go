//this is a shadowsocks server
package manager

import (
	"net"
	"sync"
	"fmt"
)




type ComChan chan int
//command for loop
const (
	NULL int = iota
	STOP
)

type server struct {
	sync.Mutex
	net.Listener
	port string //be used as id
	state int
	comChan ComChan //command channel
	local map[string]*local //1 to 1 : user.username -> local, username must be uniqueness
	format string
	started bool
}





func newServer(port string) (ss *server,err error) {

	if port == "" {
		err = newError("Cannot create a server without port")
		return
	}

	errFormat := fmt.Sprintf(serverFormat, port)
	ln, err := net.Listen("tcp", ":" + port)
	if err != nil {
		err = newError(errFormat, "create listner error:", err)
		return
	}
	ss = &server{sync.Mutex{}, ln, port, NULL, make(ComChan, serverCommand), nil, errFormat, false}
	return
}


func (self *server) addLocal(user *User) (local *local, err error) {

	local, err = newLocal(user)
	if err != nil {
		err = newError(self.format, "create local error:", err)
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

	self.stop()
	close(self.comChan)
	if err := self.Close(); err != nil {
		err = newError(self.format, "close with error:", err)
	}
	return
}

//stop the loop
func (self *server) stop() {
	if !self.isStarted() {
		go func() {
			select {
			case self.comChan <- STOP:
			}
		}()
	}
}


func (self *server) isStarted() bool {
	return self.started
}



func (self *server) start() (err error) {
	if self.started {
		err = newError(self.format, "run server error:", "has started")
	}
	go func() {
		err := self.listen()
		if err != nil {
			bd.addError(err)
		}
	}()
	return
}


func (self *server) listen() (err error) {
	self.started = true
	bd.addMsg(newLog(self.format, "start", ""))
	defer func() {
		bd.addMsg(newLog(self.format, "stop", ""))
		self.started = false
	}()
loop:
	for {
		if err != nil {
			bd.addError(err)
		}

		select{
		case com := <- self.comChan:
			if com == STOP {
				break loop
			}
		default:
		}

		conn, err := self.Accept()

		if err != nil {
			err = newError(self.format, "listener accpet error:", err)
			continue
		}
		go func() {
			err := self.handleConnect(conn)
			if err != nil {
				bd.addError(err)
			}
		}()

	}
	return
}

func(self *server) handleConnect(conn net.Conn) (err error) {

	defer func() {
		if err := recover(); err!=nil {
			if v, ok := err.(error); ok {
				conn.Write([]byte(v.Error()))
				conn.Close()
			}
		}
	}()

	user, err := getUserFormConn(conn)
	if err != nil {
		panic(newError(self.format, "get user error:", err))
	}
	//create new client and return
	local, err := self.addLocal(user)
	if err != nil {
		panic(newError(self.format, "create connect error:", err))
	}
	err = local.run()
	if err != nil {
		panic(newError(self.format, "create connect error:", err))
	}
	return
}
