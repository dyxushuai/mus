package shadowsocks

import (
	"net"
	"github.com/dropbox/godropbox/errors"
	"bytes"
	"io"
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

type Parser interface {
	Handle(c *WrapConn) (ip string, port string)
}

type FuncParser func(c *WrapConn) (ip string, port string)

func (f FuncParser) Handle(c *WrapConn) (ip string, port string) {
	return f(c)
}


func parseConn(c *WrapConn) (ip string, port string) {
	buffer := bufPool.Get()
	defer bufPool.Put(buffer)

	buffer.ReadFrom(c)

	b, err := buffer.ReadByte()

	if err != nil {
		err = errors.Wrap(err, "parese error")
		return
	}


	return
}
