package lib

import (
	"net"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"log"
)


//
type NewClientHandler interface {
	Handle(c *Client)
}

type NewClientHandlerFunc func(c *Client)

func (nc NewClientHandlerFunc) Handle(c *Client) {
	nc(c)
}


//
type ClientConnClosedHandler interface {
	Handle(c *Client, err error)
}

type ClientConnClosedHandlerFunc func(c *Client, err error)

func (nc ClientConnClosedHandlerFunc) Handle(c *Client, err error) {
	nc(c, err)
}


//
type NewDataHandler interface {
	Handle(c *Client, data []byte)
}

type NewDataHandlerFunc func(c *Client, data []byte)

func (nc NewDataHandlerFunc) Handle(c *Client, data []byte) {
	nc(c, data)
}


//the struct of config for proxy
type ProxyConfig struct {
	//listener port
	Addr string

	//encrypted string
	EncrStr string

	//encrypted method
	Method string
}

//shadowsocks server
type ProxyServer struct {
	ln net.Listener

	config 					*ProxyConfig

	// cipher
	Cip *ss.Cipher

	joins                   chan net.Conn

	onNewClientCallback     NewClientHandler
	onClientConnClosed 		ClientConnClosedHandler
	onNewData            NewDataHandler
}


func (s *ProxyServer) OnNewClient(callback NewClientHandler) {
	s.onNewClientCallback = callback
}


func (s *ProxyServer) OnClientConnClosed(callback ClientConnClosedHandler) {
	s.onClientConnClosed = callback
}


func (s *ProxyServer) OnNewData(callback NewDataHandler) {
	s.onNewData = callback
}


func (s *ProxyServer) newClient(conn net.Conn) {


	client := &Client{
		conn:   ss.NewConn(conn, s.Cip.Copy()),
		Server: s,
	}
	go client.listen()
	s.onNewClientCallback(client)
}

func (s *ProxyServer) listenChannels() {
	for {
		select {
		case conn := <-s.joins:
			s.newClient(conn)
		}
	}
}



func (s *ProxyServer) Listen() {
	go s.listenChannels()

	listener, err := net.Listen("tcp", s.config.Addr)
	if err != nil {
		log.Fatal("Error starting TCP server.")
	}
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		s.joins <- conn
	}
}



func New(conf *ProxyConfig) *ProxyServer {
	server := &ProxyServer{
		config: conf,
		joins:   make(chan net.Conn),
	}

	server.Cip = ss.NewCipher(conf.Method, conf.EncrStr)
	server.OnNewClient(func(c *Client) {})
	server.OnNewData(func(c *Client, data []byte) {})
	server.OnClientConnClosed(func(c *Client, err error) {})

	return server
}
