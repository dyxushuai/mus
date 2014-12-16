//a client connect
package manager

import (
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"github.com/cyfdecyf/leakybuf"
	"net"
	"io"
	"encoding/binary"
	"strconv"
	"time"
	"log"
	"syscall"
	"bytes"
)

type conn interface {
	net.Conn
	setTimeOut() error
}

type local struct {
	client *client
	remote *remote
}

type client struct {
	*ss.Conn
	user *User
}


type remote struct {
	net.Conn
	host string
	ip string
	port int
	extra []byte //client request content
	isHttp bool
}

//create a client and get remote
func newLocal(user *User) (local *local, err error) {
	var cipher *ss.Cipher
	cipher, err = user.cipher()
	if err != nil {
		return
	}
	local.client = &client{ss.NewConn(user.conn, cipher), user}
	local.remote, err = local.client.getRemote()
	return
}

func (self *client) setTimeOut() (err error) {
	if self.user.timeout != 0 {
		readTimeout := time.Duration(self.user.timeout) * time.Second
		err = self.SetReadDeadline(time.Now().Add(readTimeout))
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
		port int
		host string
	)
	// buf size should at least have the same size with the largest possible
	// request size (when addrType is 3, domain name has at most 256 bytes)
	// 1(addrType) + 1(lenByte) + 256(max length address) + 2(port)
	buf := make([]byte, 260)
	var n int
	// read till we get possible domain length field
	self.setTimeOut()
	if n, err = io.ReadAtLeast(self, buf, idDmLen+1); err != nil {
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
		err = newError("addr type %d not supported" + string(buf[idType]))
		return
	}
	if n < reqLen { // rare case
		self.setTimeOut()
		if _, err = io.ReadFull(self, buf[n:reqLen]); err != nil {
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
	port = int(binary.BigEndian.Uint16(buf[reqLen-2 : reqLen]))
	host = net.JoinHostPort(ip, strconv.Itoa(port))
	if err != nil {
		return
	}
	// tcp to remote
	var conn net.Conn
	conn, err = net.Dial("tcp", host)

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

	rt = &remote{conn, host, ip, port, extra, false}
	return
}




//to and from neet two goroutine
func (self *local) run() (err error) {
	closed := false

	defer func() {
		if !closed {
			self.client.Close()
			self.remote.Close()
		}
	}()

	if self.remote.isHttp {
		go func() {
			self.clientToRemote()
		}()
		self.remoteToClient()
		closed = true
	} else {
		return newError("Its not a http request")
	}
	return
}


/*
*data form remote to client
*return data length
*if data from remote to client then we will record the data length
*/
func (self *local) remoteToClient() (total int, raw_header []byte) {
	total, raw_header =  pipeThenClose(self.remote, self.client)
	self.client.user.addFlow(total)
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

var pipeBuf = leakybuf.NewLeakyBuf(nBuf, bufSize)

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
			total += size
			if err != nil {
				ss.Debug.Println("write:", err)
				break
			}
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
