//a wrapped manager for ss servers and client

package manager

import (
	"log"
	"sync"
	"net"
	"bytes"
	"syscall"
	"strings"
)
var debug ss.DebugLog

//some command
type Command int

const (
	NULL Command = iota
	WAIT Command
	OPEN Command
	CLOSE Command
)
//default commands channel with 10 maximum
var maxOfSS chan map[string]Command

type Manager struct {
	sync.Mutex
	ssServers map[string]*server //port -> ss server
}

func CreateManager() *Manager {
	return &Manager{}
}

//private method for Manager instance
func (self *Manager) hasServer(port string) bool {
	_, ok := self.ssServers[port]
	return ok
}


func (self *Manager) addClient(port string, user *User) (*client, error) {
	if self.hasServer(port) {
		return self.ssServers[port].addClient(user)
	}
	return nil, newError("no shadowsocks server listen on " + port)
}

//create a new listener with a given port
//each listener with a new goroutine
func (self *Manager) AddServerAndRun(port string) error {
	if server, err := newServer(port); err != nil {
		log.Printf("Create new server for ss at port: %s failed, err: %v\n", port, err)
		return err
	} else {
		if !self.hasServer(port) {
			self.Lock()
			self.ssServers[port] = server
			self.Unlock()
		} else {
			return newError("Application has listened this port " + port)
		}
	}
	return nil
}

//drop a existed listener
func (self *Manager) DropServer(port string) error {
	if !self.hasServer(port) {
		return newError(port + " server inexistence, drop failed")
	}
	self.Lock()
	err := self.ssServers[port].close()
	delete(self.ssServers, port)
	self.Unlock()
	return err
}


func (self *Manager) run(port string) error {
	server, err := newServer(port)
	if err != nil {
		log.Printf("Create new server for ss at port: %s failed, err: %v\n", port, err)
		return err
	}
	if !self.hasServer(port) {
		self.Lock()
		self.ssServers[port] = server
		self.Unlock()
	} else {
		return newError("Application has listened this port " + port)
	}
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Printf("accept error: %v\n", err)
			continue
		}

		user, err := getUserFormConn(conn)
		if err != nil {
			log.Printf("Error get passeoed: %s %v\n", port, err)
			conn.Close()
			continue
		}

		if user.currentFlow < user.limit {
			log.Printf("Error user runover: %s\n", user.limit)
			conn.Close()
			continue
		}
		client, err := self.addConn(port, user)
		if err != nil {
			log.Printf("Add client failed: %s %v\n", port, err)
			conn.Close()
			continue
		}
		go handleConnection(client)
	}
}

func getRequest(conn *ss.Conn) (host string, extra []byte, err error) {
	const (
		idType  = 0 // address type index
		idIP0   = 1 // ip addres start index
		idDmLen = 1 // domain address length index
		idDm0   = 2 // domain address start index

		typeIPv4 = 1 // type is ipv4 address
		typeDm   = 3 // type is domain address
		typeIPv6 = 4 // type is ipv6 address

		lenIPv4   = 1 + net.IPv4len + 2 // 1addrType + ipv4 + 2port
		lenIPv6   = 1 + net.IPv6len + 2 // 1addrType + ipv6 + 2port
		lenDmBase = 1 + 1 + 2           // 1addrType + 1addrLen + 2port, plus addrLen
	)

	// buf size should at least have the same size with the largest possible
	// request size (when addrType is 3, domain name has at most 256 bytes)
	// 1(addrType) + 1(lenByte) + 256(max length address) + 2(port)
	buf := make([]byte, 260)
	var n int
	// read till we get possible domain length field
	SetReadTimeout(conn)
	if n, err = io.ReadAtLeast(conn, buf, idDmLen+1); err != nil {
		return
	}

	reqLen := -1
	switch buf[idType] {
	case typeIPv4:
		reqLen = lenIPv4
	case typeIPv6:
		reqLen = lenIPv6
	case typeDm:
		reqLen = int(buf[idDmLen]) + lenDmBase
	default:
		err = errors.New(fmt.Sprintf("addr type %d not supported", buf[idType]))
		return
	}

	if n < reqLen { // rare case
		SetReadTimeout(conn)
		if _, err = io.ReadFull(conn, buf[n:reqLen]); err != nil {
			return
		}
	} else if n > reqLen {
		// it's possible to read more than just the request head
		extra = buf[reqLen:n]
	}

	// Return string for typeIP is not most efficient, but browsers (Chrome,
	// Safari, Firefox) all seems using typeDm exclusively. So this is not a
	// big problem.
	switch buf[idType] {
	case typeIPv4:
		host = net.IP(buf[idIP0 : idIP0+net.IPv4len]).String()
	case typeIPv6:
		host = net.IP(buf[idIP0 : idIP0+net.IPv6len]).String()
	case typeDm:
		host = string(buf[idDm0 : idDm0+buf[idDmLen]])
	}
	// parse port
	port := binary.BigEndian.Uint16(buf[reqLen-2 : reqLen])
	host = net.JoinHostPort(host, strconv.Itoa(int(port)))
	return
}

func handleConnection(client *client) {
	var host string
	var size = 0
	var raw_req_header, raw_res_header []byte
	var is_http = false
	var res_size = 0
	var req_chan = make(chan []byte)


	// function arguments are always evaluated, so surround debug statement
	// with if statement
	if debug {
		debug.Printf("new client %s->%s\n", client.conn.RemoteAddr().String(), client.conn.LocalAddr())
	}
	closed := false
	defer func() {
		if debug {
			debug.Printf("closed pipe %s<->%s\n", client.conn.RemoteAddr(), host)
		}
		if !closed {
			client.conn.Close()
		}
	}()

	host, extra, err := getRequest(client.conn)
	if err != nil {
		log.Println("error getting request", client.conn.RemoteAddr(), client.conn.LocalAddr(), err)
		return
	}
	debug.Println("connecting", host)
	remote, err := net.Dial("tcp", host)
	if err != nil {
		if ne, ok := err.(*net.OpError); ok && (ne.Err == syscall.EMFILE || ne.Err == syscall.ENFILE) {
			// log too many open file error
			// EMFILE is process reaches open file limits, ENFILE is system limit
			log.Println("dial error:", err)
		} else {
			log.Println("error connecting to:", host, err)
		}
		return
	}
	defer func() {
		if is_http{
			tmp_req_header := <-req_chan
			buffer := bytes.NewBuffer(raw_req_header)
			buffer.Write(tmp_req_header)
			raw_req_header = buffer.Bytes()
		}
		close(req_chan)
		if !closed {
			remote.Close()
		}
	}()
	// write extra bytes read from

	is_http, extra, _ = checkHttp(extra, client.conn)
	if strings.HasSuffix(host, ":80") {
		is_http = true
	}
	raw_req_header = extra
	res_size, err = remote.Write(extra)
//	storage.IncrSize("flow:" + user.Name, res_size)
//	storage.ZincrbySize("flow:" + user.Name, host, res_size)
	size += res_size
	if err != nil {
		debug.Println("write request extra error:", err)
		return
	}

	if debug {
		debug.Printf("piping %s<->%s", client.conn.RemoteAddr(), host)
	}

	go func() {
		_, raw_header := PipeThenClose(conn, remote, ss.SET_TIMEOUT, is_http, false, host, user)
		if is_http {
			req_chan<-raw_header
		}
	}()

	res_size, raw_res_header = PipeThenClose(remote, conn, ss.NO_TIMEOUT, is_http, true, host, user)
	size += res_size
	closed = true
	return
}


func checkHttp(extra []byte, conn *ss.Conn) (is_http bool, data []byte, err error) {
	var buf []byte
	var methods = []string{"GET", "HEAD", "POST", "PUT", "TRACE", "OPTIONS", "DELETE"}
	is_http = false
	if extra == nil || len(extra) < 10 {
		buf = make([]byte, 10)
		if _, err = io.ReadFull(conn, buf); err != nil {
			return
		}
	}

	if buf == nil {
		data = extra
	} else if extra == nil {
		data = buf
	}else {
		buffer := bytes.NewBuffer(extra)
		buffer.Write(buf)
		data = buffer.Bytes()
	}

	for _, method := range methods {
		if bytes.HasPrefix(data, []byte(method)) {
			is_http = true
			break
		}
	}
	return
}
