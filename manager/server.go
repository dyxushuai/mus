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



type Server struct {
	mu sync.Mutex

	Port          string       `json:"port"`
	Method        string       `json:"method"`
	Password      string       `json:"password"`
	Current       int64        `json:"current"`
	Limit         int64        `json:"limit"`
	Timeout       int64        `json:"timeout"`

	listener      net.Listener
	comChan       ComChan          //command channel
	local		  map[string]*local //1 to 1 : remote addr -> local
	format        string
	started       bool
	cipher        *ss.Cipher
	store         *Storage
}

func newServer(port, method, password string, limit, timeout int64, redis *Storage) (server *Server,err error) {

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

	server = &Server{
		Port: port,
		Method: method,
		Password: password,
		Timeout: timeout,
		Current: 0,
		Limit: limit,
		listener: ln,
		comChan: make(chan ComChan),
		local: make(map[string]*local),
		format: errFormat,
		started: false,
		cipher: cipher,
		store: redis,
	}

	return
}

func (self *Server) addLocal(conn net.Conn) (local *local, err error) {

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
func (self *Server) isStarted() bool {
	return self.started
}
func (self *Server) isOverFlow() bool {
	return self.Current > self.Limit
}

func (self *Server) Close() (err error) {
	//first stop the loop
	//second close the chan
	//third close the listener

	self.stop()
	close(self.comChan)
	if err := self.listener.Close(); err != nil {
		err = newError(self.format, "close with error:", err)
	}
	return
}

func (self *Server) Start() (err error) {
	if self.started {
		err = newError(self.format, "run server error:", "has started")
	}
	go func() {
		err := self.listen()
		if err != nil {
			return
		}
	}()
	return
}

//stop the loop
func (self *Server) Stop() {
	if !self.isStarted() {
		go func() {
			select {
			case self.comChan <- STOP:
			}
		}()
	}
}

func (self *Server) listen() (err error) {
	self.started = true
	log.Info("server at port: %s started", self.Port)
	defer func() {
		self.started = false
		log.Info("server at port: %s stoped", self.Port)
	}()
loop:
	for {

		select{
		case com := <- self.comChan:
			if com == STOP {
				break loop
			}
		default:
		}

		conn, err := self.listener.Accept()
		if err != nil {
			err = newError(self.format, "listener accpet error:", err)
			log.Debug(err)
			continue
		}
		go func() {
			flow, err := self.handleConnect(conn)
			//TODO: use `flow`
			_ = flow
			if err != nil {
				log.Debug(err)
			}
		}()

	}
	return
}

func (self *Server) handleConnect(conn net.Conn) (flow int, err error) {

	defer func () {
		conn.Close()
	}()

	local, err := self.addLocal(conn)

	if err != nil {
		return
	}
	flow, err = local.run()
	if err != nil {
		return
	}
	return
}


