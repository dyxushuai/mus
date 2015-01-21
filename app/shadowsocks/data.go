package shadowsocks

import (
	"github.com/JohnSmithX/mus/app/shadowsocks/lib"
	"net"
	"io"
)





type Data struct {

}

func (d *Data) Parse() {

}

func (d *Data) Handle(*lib.Client, data []byte) {

}


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
		err = utils.NewError(self.father.format, "read domain form client error:", err)
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
		err = utils.NewError(self.father.format, "type of request address is not supported", string(buf[idType]))
		return
	}
	if n < reqLen { // rare case
		self.setTimeOut()
		if _, err = io.ReadFull(self, buf[n:reqLen]); err != nil {
			err = utils.NewError(self.father.format, "read addr form client error:", err)
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
			err = utils.NewError(self.father.format, "dial error:", err)
		} else {
			err = utils.NewError(self.father.format, "error connecting to: " + host, err)
		}
		return
	}

	rt = &remote{conn, host, ip, port, extra, false, self.father}
	err = rt.checkMethod()
	return
}


