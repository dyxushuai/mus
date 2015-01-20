package shadowsocks

import (
	"net"
	"github.com/dropbox/godropbox/errors"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"bytes"
)


var (
	//buffer pool
	bufPool *BufferPool = NewBufferPool(2048)
)


//shadowsocks server
type ProxyServer struct {
	net.Listener

	//listener port
	Port string
	//encrypted string
	EncrStr string
	//encrypted method
	Method string


	//the num of activated connection
	ConnCount *counter

	// cipher
	Cip *ss.Cipher
	//remote connection maker
	Dial func(network, addr string) (net.Conn, error)
}


func (ps *ProxyServer) Accept() (c net.Conn, err error) {
	c, err = ps.Listener.Accept()
	if err != nil {
		err = errors.Wrapf(err, "at %s new connection error", ps.Port)
		return
	}

	// Wrap the returned connection so we're able to observe
	// when it is closed
	c = &WrapConn{Conn: c, ConnCount: ps.ConnCount}

	err = ps.CipMaker(c)
	if err != nil {
		return
	}

	// Count it
	ps.ConnCount.Lock()
	ps.ConnCount.c++
	ps.ConnCount.Unlock()
	return
}

func (ps *ProxyServer) CipMaker(wc *WrapConn) (err error) {
	if ps.Cip == nil {
		ps.Cip, err = ss.NewCipher(ps.Method, ps.EncrStr)
		if err != nil {
			err = errors.Wrapf(err, "at %s new cipher error", ps.Port)
			return
		}
	}
	if wc.Cipher == nil {
		wc.Cipher = ps.Cip.Copy()
	}
	return
}

