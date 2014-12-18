//users info

package models

import (
	"io"
	"net"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
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


func getUserFormConn(conn net.Conn) (user *User, err error) {
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
		err = newError("version %s not supported" + string(buf[idVersion]))
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


