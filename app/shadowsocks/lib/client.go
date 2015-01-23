package lib


import (
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"github.com/dropbox/godropbox/errors"
	"github.com/oxtoacart/bpool"
	"net"
	"io"
	"encoding/binary"
	"strconv"
	"syscall"
)

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
	bytePool = bpool.NewBytePool(4096, 2048)
)


type SSClienter interface {
	net.Conn
	Remote() net.Conn
}

type client struct {
	*ss.Conn
	server 		*ProxyServer
	remote		*remote
	closed		bool
}

type remote struct {
	net.Conn
	closed		bool
}

func (r *remote) Close() (err error) {
	if !r.closed {
		err = r.Conn.Close()
		r.closed = true
	}
	return
}

func (r *remote) Write(b []byte) (n int, err error) {
	if !r.closed {
		n, err = r.Conn.Write(b)
	}
	return
}

func (r *remote) Read(b []byte) (n int, err error) {
	if !r.closed {
		n, err = r.Conn.Read(b)
	}
	return
}


func (c *client) Close() (err error) {
	if !c.closed {
		err = c.Conn.Conn.Close()
		c.closed = true
	}
	return
}

func (c *client) Write(b []byte) (n int, err error) {
	if !c.closed {
		n, err = c.Conn.Write(b)
	}
	return
}

func (c *client) Read(b []byte) (n int, err error) {
	if !c.closed {
		n, err = c.Conn.Read(b)
	}
	return
}

func (c *client) rListen() {
	defer c.remote.Close()

	buf := bytePool.Get()
	defer bytePool.Put(buf)

	flow := 0
	defer func() {
		c.server.CallbackMethods.Record(flow)
	}()

	for {
		n, err := c.remote.Read(buf)

		if n > 0 {
			flow += n
			e := c.server.CallbackMethods.RemoteNewData(c, buf[:n])
			if e != nil {
				return
			}
		}

		if err != nil {
			c.server.CallbackMethods.RemoteReadErr(c, err)
			return
		}
	}
}

func (c *client) newRemote(conn net.Conn) {

	c.remote = &remote{
		Conn: conn,
	}
	go c.rListen()
	c.server.CallbackMethods.NewRemote(c)
}

func (c *client) parse() (conn net.Conn, extra []byte, err error) {

	var (
		ip string
		port string
		host string
	)
	// buf size should at least have the same size with the largest possible
	// request size (when addrType is 3, domain name has at most 256 bytes)
	// 1(addrType) + 1(lenByte) + 256(max length address) + 2(port)
	buf := make([]byte, 260)

	// read till we get possible domain length field

	n, err := io.ReadAtLeast(c.Conn, buf, idDmLen+1)

	if err != nil {
		err = errors.Newf("read domain form client error: %v", err)
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
		err = errors.Newf("type of request (%s) address is not supported", string(buf[idType]))
		return
	}
	if n < reqLen { // rare case
		if _, err = io.ReadFull(c.Conn, buf[n:reqLen]); err != nil {
			err = errors.Newf("read addr form client error: %v", err)
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
	conn, err = net.Dial("tcp", host)

	if err != nil {
		if ne, ok := err.(*net.OpError); ok && (ne.Err == syscall.EMFILE || ne.Err == syscall.ENFILE) {
			// log too many open file error
			// EMFILE is process reaches open file limits, ENFILE is system limit
			err = errors.Newf("get remote dial error: %v", err)
		} else {
			err = errors.Newf("error connecting to remote: %s %v",  host, err)
		}
		return
	}

	return
}


func (c *client) listen() {

	conn, extra, er := c.parse()

	if  er != nil {
		return
	}
	
	c.newRemote(conn)

	flow := 0
	defer func() {
		c.server.CallbackMethods.Record(flow)
	}()


	if extra != nil && len(extra) != 0 {
		c.server.CallbackMethods.ClientNewData(c, extra )
	}
	
	buf := bytePool.Get()
	defer bytePool.Put(buf)

	for {
		n, err := c.Read(buf)

		if n > 0 {
			flow += n
			e := c.server.CallbackMethods.ClientNewData(c, buf[:n])
			if e != nil {
				return
			}

		}

		if err != nil {
			c.server.CallbackMethods.ClientReadErr(c, err)
			return
		}

	}
}


func (c *client) Remote() net.Conn {
	return c.remote
}
