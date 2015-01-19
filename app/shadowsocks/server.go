//this is a shadowsocks server
package shadowsocks

import (
	"net"
	"sync"
	"fmt"
	"strings"
	"github.com/JohnSmithX/mus/app/utils"
	"github.com/JohnSmithX/mus/app/db"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
)


//server config
const (
	//format output
	serverFormat string = "proxy server at port %s : %%s %%v"
	//command for loop
	NULL int = iota
	STOP

)


type (
	ComChan chan int
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
	
	recorder		db.IStorage
	listener      	net.Listener
	comChan       	ComChan          	//command channel
	local		  	map[string]*local //1 to 1 : remote addr -> local
	format        	string
	cipher        	*ss.Cipher
}

func NewServer(port, method, password string, limit, timeout int64 ,recorder db.IStorage) (server *Server,err error) {
	if port == "" {
		err = utils.NewError("Cannot create a server without port")
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

	err = server.InitServer(recorder)
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
		err = utils.NewError(self.format, "create local error:", err)
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


func (self *Server) addFlow(flow int) (err error) {

	self.doWithLock(func() {
		_, err = self.recorder.IncrSize("flow:" + self.Port, flow)
		self.Current += int64(flow)
		utils.Debug(err)
	})
	return
}

func (self *Server) listen() {
	self.Started = true
	utils.Info("server at port: %s started", self.Port)
	defer func() {
		self.Started = false
		utils.Info("server at port: %s stoped", self.Port)
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
			err = utils.NewError(self.format, "listener accpet error:", err)
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
func (self *Server) InitServer(recorder db.IStorage) (err error) {

	self.comChan = make(ComChan)
	self.local =  make(map[string]*local)
	self.Started = false
	self.recorder = recorder

	errFormat := fmt.Sprintf(serverFormat, self.Port)

	ln, err := net.Listen("tcp", ":" + self.Port)
	if err != nil {
		err = utils.NewError(errFormat, "create listner error:", err)
		return
	}

	cipher, err := ss.NewCipher(self.Method, self.Password)
	if err != nil {
		err = utils.NewError(errFormat, "create cipher error:", err)
		return
	}

	self.format = errFormat
	self.listener = ln
	self.cipher = cipher

	return
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
		err = utils.NewError(self.format, "run server error:", "has started")
		return
	}

	go func () {
		self.listen()
	}()
	return
}

func (self *Server) Stop() (err error) {
	if !self.isStarted() {
		err = utils.NewError(self.format, "run server error:", "has stopped")
		return
	}
	go func() {
		select {
		case self.comChan <- STOP:
		}
	}()
	return
}

func (self *Server) Destroy() (err error) {
	//first stop the loop
	//second close the chan
	//third close the listener

	self.Stop()
	close(self.comChan)
	if err := self.listener.Close(); err != nil {
		err = utils.NewError(self.format, "close with error:", err)
	}
	return
}

func (self *Server) Logs() (result string, err error) {return}

func (self *Server) Flow() (result string, err error) {return}

func (self *Server) Key() string {
	return self.Port
}
