//this is a shadowsocks server
package shadowsocks

import (
	"net"
	"sync"
	"fmt"
	"strings"
	"encoding/json"
	"github.com/JohnSmithX/mus/app/utils"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
)


//server config
const (
	//format output
	serverFormat string = "proxy server at port %s : %%s %%v"
)


var Log utils.Verbose
var newError = utils.NewError


type ComChan chan int

type Recorder interface {
	IncrSize(string, int) (int64, error)
}

//command for loop
const (
	NULL int = iota
	STOP
)


type Server struct {
	mu sync.Mutex

	Port 			string			`json:"port"`
	Method       	string       	`json:"method"`
	Password      	string       	`json:"password"`
	Limit         	int64        	`json:"limit"`
	Timeout       	int64        	`json:"timeout"`
	Current       	int64       	`json:"current"`
	Started			bool			`json:"started"`// the state of server
	
	recorder		Recorder
	listener      	net.Listener
	comChan       	ComChan          	//command channel
	local		  	map[string]*local //1 to 1 : remote addr -> local
	format        	string
	cipher        	*ss.Cipher
}


func NewServer(port, method, password string, limit, timeout int64 ,recorder Recorder) (server *Server,err error) {
	if port == "" {
		err = newError("Cannot create a server without port")
		return
	}

	server = &Server{
		Port: port,
		Method: method,
		Password: password,
		Timeout: timeout,
		Limit: limit,
		Current: 0,
		recorder: recorder,
	}

	err = server.initServer()
	return
}

func (self *Server) initServer() (err error) {
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
	self.comChan = make(ComChan)
	self.local =  make(map[string]*local)
	self.Started = false
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
	return self.Started
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

func (self *Server) addFlow(flow int) (err error) {
	_, err = self.recorder.IncrSize("flow:" + self.Port, flow)
	self.current += int64(flow)
	utils.Debug(err)
	return
}

func (self *Server) listen() {
	self.Started = true
	Log.Info("server at port: %s started", self.Port)
	defer func() {
		self.Started = false
		Log.Info("server at port: %s stoped", self.Port)
	}()
loop:
	for {

		if self.isOverFlow() {
			break loop
		}

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
			utils.Debug(err)
			continue
		}
		go self.handleConnect(conn)

	}
	return
}

func (self *Server) handleConnect(conn net.Conn) (flow int, err error) {

	defer func () {
		utils.Debug(err)
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

	err = self.addFlow(flow)
	return
}

//interface
func (self *Server) JSON() string {
	data, _ := json.Marshal(self)
	return string(data)
}

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

func (self *Server) Logs() {}

func (self *Server) Flow() {}

