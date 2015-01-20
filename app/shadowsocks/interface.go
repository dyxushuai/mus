package  shadowsocks

import (
	"github.com/JohnSmithX/mus/app/db"
	"net"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
)

type ShadowsocksServer interface {
	//initialize function with db to recorder flow
	InitServer(db.IStorage) error

	//hold
	Logs() (string, error)

	Flow() (string, error)
	//actions
	ReStart() error

 	Start() error

	Stop() error

	Destroy() error

	Key() string
}



type ReqHandler interface {
	Handle(cn net.Conn, EncrStr string, Method string) ss.Conn
}
