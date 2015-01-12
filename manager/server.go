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
	manager       *Manager
}


func newServer(port, method, password string, limit, timeout int64, father *Manager) (server *Server,err error) {
	if port == "" {
		err = newError("Cannot create a server without port")
		return
	}

	server = &Server{
		Port: port,
		Method: method,
		Password: password,
		Timeout: timeout,
		Current: 0,
		Limit: limit,
		comChan: make(ComChan),
		local: make(map[string]*local),
		started: false,
	}

	err = server.initServer(father)
	return
}

func (self *Server) initServer(father *Manager) (err error) {
	errFormat := fmt.Sprintf(serverFormat, self.Port)
	ln, err := net.Listen("tcp", ":" + self.Port)
	if err != nil {
		err = newError(errFormat, "create listner error:", err)
		return
	}

	cipher, err := ss.NewCipher(self.Method, self.Password)
	if err != nil {
		err = newError(errFormat, "create cipher error:", err)
		return
	}

	self.format = errFormat
	self.listener = ln
	self.cipher = cipher
	self.manager = father
	return
}


func (self *Server) doWithLock(fn func()) {
	self.mu.Lock()
	defer self.mu.Unlock()
	fn()
}

func (self *Server) addLocal(conn net.Conn) (local *local, err error) {

	cipher := self.cipher.Copy()
	ssconn := ss.NewConn(conn, cipher)

	local, err = newLocal(self, ssconn)
	if err != nil {
		err = newError(self.format, "create local error:", err)
		return
	}

	ip := strings.Split(conn.RemoteAddr().String(), ":")[0]
	self.doWithLock(func() {
		self.local[ip] = local
	})

	return
}

func (self *Server) isStarted() bool {
	return self.started
}

func (self *Server) isOverFlow() bool {
	return self.Current > self.Limit
}

func (self *Server) destroy() (err error) {
	//first stop the loop
	//second close the chan
	//third close the listener

	self.Stop()
	close(self.comChan)
	if err := self.listener.Close(); err != nil {
		err = newError(self.format, "close with error:", err)
	}
	return
}

//record flow
func (self *Server) recFlow(flow int) (err error) {
	_, err = self.manager.store.IncrSize(flowPrefix + self.Port, flow)
	return
}

func (self *Server) listen() {
	self.started = true
	Log.Info("server at port: %s started", self.Port)
	defer func() {
		self.started = false
		Log.Info("server at port: %s stoped", self.Port)
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
			Log.Debug(err.Error())
			continue
		}
		go func() {
			flow, er := self.handleConnect(conn)
			//TODO: use `flow`
			if er != nil {
				Log.Debug(er.Error())
				return
			}
			er = self.recFlow(flow)
			Log.Debug(er.Error())
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



//TODO: API
//POST /api/servers/:id/restart
func (self *Server) ReStart() (err error) {
	if self.isStarted() {
		err = self.Stop()
		if err != nil {
			return
		}
	}
	err = self.Start()
	return
}

//POST /api/servers/:id/start
func (self *Server) Start() (err error) {
	if self.isStarted() {
		err = newError(self.format, "run server error:", "has started")
		return
	}
	go func () {
		self.listen()
	}()
	return
}

//POST /api/servers/:id/stop
func (self *Server) Stop() (err error) {
	if !self.isStarted() {
		err = newError(self.format, "run server error:", "has stopped")
		return
	}
	go func() {
		select {
		case self.comChan <- STOP:
		}
	}()
	return
}
