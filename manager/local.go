//a client connect
package manager

import (
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"net"
	"io"
	"encoding/binary"
	"strconv"
	"time"
	"syscall"
	"bytes"
	"fmt"
)


type conn interface {
	net.Conn
	setTimeOut() error
}

type local struct {
	client *client
	remote *remote
	format string
	father *Server
}

type client struct {
	*ss.Conn
	father *local
}


type remote struct {
	net.Conn
	host string
	ip string
	port string
	extra []byte //client request content
	isHttp bool
	father *local

}

//create a client and get remote
func newLocal(sserver *Server,conn *ss.Conn) (l *local, err error) {

	format := fmt.Sprintf(localFormat, conn.RemoteAddr())

	l = new(local)

	l.format = format
	l.father = sserver
	l.client = &client{conn, l}

	l.remote, err = l.client.getRemote()

	return
}

func (self *client) setTimeOut() (err error) {
	if self.father.father.Timeout != 0 {
		readTimeout := time.Duration(self.father.father.Timeout) * time.Second
		err = self.SetReadDeadline(time.Now().Add(readTimeout))
		err = newError(self.father.format, "client set timeout error:", err)
	}
	return
}

func (self *remote) setTimeOut() (err error) {
	return
}


/*
*filter http request from client
*exclude other protocol
*/
func (self *remote) checkMethod() (err error) {
	var methods = []string{"GET", "HEAD", "POST", "PUT", "TRACE", "OPTIONS", "DELETE"}
	self.isHttp = false

	//
	if self.extra == nil {
		err = newError("No request content form client")
		return
	}

	for _, method := range methods {
		if bytes.HasPrefix(self.extra, []byte(method)) {
			self.isHttp = true
			break
		}
	}
	return
}
//get remote info
func (self *client) getRemote() (rt *remote, err error) {
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
	var (
		extra []byte
		ip string
		port string
		host string
	)
	// buf size should at least have the same size with the largest possible
	// request size (when addrType is 3, domain name has at most 256 bytes)
	// 1(addrType) + 1(lenByte) + 256(max length address) + 2(port)
	buf := make([]byte, 260)

	// read till we get possible domain length field
	self.setTimeOut()
	n, err := io.ReadAtLeast(self, buf, idDmLen+1)

	if err != nil {
		err = newError(self.father.format, "read domain form client error:", err)
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
		err = newError(self.father.format, "type of request address is not supported", string(buf[idType]))
		return
	}
	if n < reqLen { // rare case
		self.setTimeOut()
		if _, err = io.ReadFull(self, buf[n:reqLen]); err != nil {
			err = newError(self.father.format, "read addr form client error:", err)
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
		ip = net.IP(buf[idIP0 : idIP0+net.IPv4len]).String()
	case typeIPv6:
		ip = net.IP(buf[idIP0 : idIP0+net.IPv6len]).String()
	case typeDm:
		ip = string(buf[idDm0 : idDm0+buf[idDmLen]])
	}
	// parse port
	port = strconv.Itoa(int(binary.BigEndian.Uint16(buf[reqLen-2 : reqLen])))
	host = net.JoinHostPort(ip, port)
	// tcp to remote
	conn, err := net.Dial("tcp", host)

	if err != nil {
		if ne, ok := err.(*net.OpError); ok && (ne.Err == syscall.EMFILE || ne.Err == syscall.ENFILE) {
			// log too many open file error
			// EMFILE is process reaches open file limits, ENFILE is system limit
			err = newError(self.father.format, "dial error:", err)
		} else {
			err = newError(self.father.format, "error connecting to: " + host, err)
		}
		return
	}

	rt = &remote{conn, host, ip, port, extra, false, self.father}
	err = rt.checkMethod()
	return
}




//to and from neet two goroutine
func (self *local) run() (flow int, err error) {
	closed := false

	defer func() {
		if !closed {
			self.client.Close()
		}
	}()


	go func() {
		self.clientToRemote()
	}()
	flow, _ = self.remoteToClient()
	closed = true
//	if self.remote.isHttp {
//		go func() {
//			self.clientToRemote()
//		}()
//		flow, _ = self.remoteToClient()
//		closed = true
//	} else {
//		err = newError(self.format, "request a non-http method", "")
//	}
	return
}


/*
*data form remote to client
*return data length
*if data from remote to client then we will record the data length
*/
func (self *local) remoteToClient() (total int, raw_header []byte) {
	total, raw_header =  pipeThenClose(self.remote, self.client)
//	self.client.user.addFlow(total)
	return
}



func (self *local) clientToRemote() (total int, raw_header []byte) {
	total, raw_header =  pipeThenClose(self.client, self.remote)
	return
}


const (
	bufSize = 4096
	nBuf = 2048
)

var pipeBuf = ss.NewLeakyBuf(nBuf, bufSize)

//from source to  destination ->
func pipeThenClose(src, dst conn) (total int, raw_header []byte) {
	defer dst.Close()

	buf := pipeBuf.Get()
	defer pipeBuf.Put(buf)
	var buffer = bytes.NewBuffer(nil)
	var is_end = false
	var size int

	for {
		src.setTimeOut()
		n, err := src.Read(buf)
		log.Info(string(buf))
		// read may return EOF with n > 0
		// should always process n > 0 bytes before handling error
		if n > 0 {
			if  !is_end {
				buffer.Write(buf)
				raw_header = buffer.Bytes()
				lines := bytes.SplitN(raw_header, []byte("\r\n\r\n"), 2)
				if len(lines) == 2 {
					is_end = true
				}
			}
			size, err = dst.Write(buf[0:n])
			if err != nil {
				ss.Debug.Println("write:", err)
				break
			}
			total += size
		}
		if err != nil || n == 0 {
			//==
			// Always "use of closed network connection", but no easy way to
			// identify this specific error. So just leave the error along for now.
			// More info here: https://code.google.com/p/go/issues/detail?id=4373
			break
		}
	}
	return
}
