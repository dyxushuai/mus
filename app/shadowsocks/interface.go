package  shadowsocks

import (
	"github.com/JohnSmithX/mus/app/db"
)

type ShadowsocksServer interface {
	//initialize function with db to recorder flow
	InitServer(db.IStorage) error

	//json text
	JSON() (string, error)

	Logs() (string, error)

	Flow() (string, error)
	//actions
	ReStart() error

 	Start() error

	Stop() error

	Destroy() error

	Key() string
}
