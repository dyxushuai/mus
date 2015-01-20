package shadowsocks

import (
	"sync"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
)


type WrapConn struct {
	ss.Conn
	ConnCount *counter
}

func (wc *WrapConn) Close() error {
	wc.ConnCount.Lock()
	wc.ConnCount.c--
	wc.ConnCount.Unlock()
	return wc.Conn.Close()
}


type counter struct {
	sync.Mutex
	c int
}

func (c *counter) Get() int {
	c.Lock()
	defer c.Unlock()
	return c.c
}

