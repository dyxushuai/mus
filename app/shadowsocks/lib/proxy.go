package lib

import (
	"net"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"log"
	"time"
	"sync"
	"github.com/dropbox/godropbox/errors"
)


type CallbackInterface interface {
	//client event
	NewClient(c SSClienter)

	ClientReadErr(c SSClienter, err error)

	ClientNewData(c SSClienter, data []byte) error
	//remote event
	NewRemote(c SSClienter)

	RemoteReadErr(c SSClienter, err error)

	RemoteNewData(c SSClienter, data []byte) error

	Record(i int)
}



//the struct of config for proxy
type ProxyConfig struct {
	//listener port
	Addr string

	//encrypted string
	EncrStr string

	//encrypted method
	Method string

	//read and write timeout
	//second
	Timeout time.Duration
}

//shadowsocks server
type ProxyServer struct {
	mu						sync.Mutex

	ln 						net.Listener

	config 					*ProxyConfig

	// cipher
	Cip 					*ss.Cipher

	stop      				chan bool
	Stopped   				bool

	joins                   chan net.Conn

	CallbackMethods			CallbackInterface
}



func (s *ProxyServer) doWithLock(fn func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fn()
}

func (s *ProxyServer) newClient(conn net.Conn) {

	conn.SetReadDeadline(time.Now().Add(s.config.Timeout * time.Second))

	client := &client{
		Conn:   ss.NewConn(conn, s.Cip.Copy()),
		server: s,
	}

	go client.listen()
	s.CallbackMethods.NewClient(client)
}

func (s *ProxyServer) listenChannels() {
	defer func() {
		recover()
	}()
	for {
		select {
		case conn := <-s.joins:
			s.newClient(conn)
		case b := <-s.stop:
			if b {
				return
			}
		}
	}
}


func (s *ProxyServer) Stop() {

	s.stop <-true
}

func (s *ProxyServer) Listen() {
	go s.listenChannels()

	s.doWithLock(func() {
		s.Stopped = false
	})


	listener, err := net.Listen("tcp", s.config.Addr)
	if err != nil {
		log.Printf("Error starting TCP server. %v", err)
		return
	}
	defer listener.Close()

	for {
		select {
		case b := <-s.stop:
			if b {
				if !s.Stopped {
					s.doWithLock(func() {
						s.Stopped = true
					})
					return
				}

			}
		default:
		}
		conn, _ := listener.Accept()
		s.joins <-conn

	}
}

func (s *ProxyServer) IsStopped() (b bool) {
	s.doWithLock(func() {
		b = s.Stopped
	})
	return
}

func (s *ProxyServer) SetCallbacks(cb CallbackInterface) {
	s.CallbackMethods = cb
}

func New(conf *ProxyConfig) (ps *ProxyServer, err error) {
	ps = &ProxyServer{
		config: conf,
		joins:   make(chan net.Conn),
		stop:    make(chan bool),
		Stopped: true,
	}
	ps.Cip, err = ss.NewCipher(conf.Method, conf.EncrStr)
	if err != nil {
		err = errors.Newf("create cipher error: %v", err)
	}
	return
}
