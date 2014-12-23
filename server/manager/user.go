//users info

package manager

import (
	"net"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"io"
)


type User struct {
	username      string
	password      string
	method        string
	timeout 	  int
	currentFlow   float64
	limit		  float64
	conn          net.Conn
}

//Encrypted string
func (self *User) cipher() (*ss.Cipher, error) {
	return ss.NewCipher(self.method, self.password)
}

func (self *User) addFlow(size int) (err error) {
	return nil
}

//func getUserFormConn(conn net.Conn) (u *User, err error)  {
//	u = &User{}
//	u.username = "xus"
//	u.password = "123"
//	u.method = "table"
//	u.timeout = 600
//	u.currentFlow = 15000000
//	u.limit = 20000000
//	u.conn = conn
//	return
//
//}

func getUserFormConn(conn net.Conn) (user *User, err error) {
	//get user form db
	//check if overflow
	const (
		idVersion = 0 // \x01
		version   = 1
		idUser    = 1 // \x02
	)
	buf := make([]byte, 2)
	if _, err = io.ReadFull(conn, buf); err != nil {
		return
	}
	switch buf[idVersion] {
	case version:
		break
	default:
		err = newError("version %s not supported ", string(buf[idVersion]))
		return
	}
	userLen := buf[idUser]
	username := make([]byte, userLen)
	if _, err = io.ReadFull(conn, username); err != nil {
		return
	}
//	user, err = storage.Get("user:" + string(username))

	return

}

