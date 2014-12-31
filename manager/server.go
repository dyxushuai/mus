//this is a shadowsocks server
package manager

import (
	"net"
	"sync"
	"fmt"
	"strings"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
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
	comChan ComChan //command channel
	local map[string]*local //1 to 1 : remote addr -> local
	format string
	started bool
	method string //encrypt method
	password string //secret key
	cipher *ss.Cipher
	timeout int64
}




func newServer(port, method, password string, timeout int64) (sserver *server,err error) {

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

	cipher, err := ss.NewCipher(method, password)
	if err != nil {
		err = newError(errFormat, "create cipher error:", err)
		return
	}

	sserver = &server{sync.Mutex{}, ln, port, make(ComChan, serverCommand), make(map[string]*local), errFormat, false, method, password, cipher, timeout}
	return
}


func (self *server) addLocal(conn net.Conn) (local *local, err error) {

	cipher := self.cipher.Copy()
	ssconn := ss.NewConn(conn, cipher)

	local, err = newLocal(self, ssconn)
	if err != nil {
		err = newError(self.format, "create local error:", err)
		return
	}

	self.Lock()
	defer func () {
		self.Unlock()
	}()
	ip := strings.Split(conn.RemoteAddr().String(), ":")[0]
	self.local[ip] = local

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

	defer func () {
		conn.Close()
	}()

	local, err := self.addLocal(conn)

	if err != nil {
		return
	}
	err = local.run()
	if err != nil {
		return
	}
	return
}
