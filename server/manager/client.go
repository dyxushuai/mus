//a client connect

package manager

import (
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"

)

type client struct {
	user *User
	conn *ss.Conn
}

func newClient(user *User) (*client, error) {
	if cipher, err := user.cipher(); err != nil {
		return nil, err
	} else {
		ssConn := ss.NewConn(user.conn, cipher)
		return &client{user: user, conn: ssConn }, nil
	}
}
