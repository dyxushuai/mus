package  shadowsocks

type ShadowsocksServer interface {
	//initialize function
	InitServer() error

	//json text
	JSON() (string, error)

	Logs() (string, error)

	Flow() (string, error)
	//actions
	ReStart() error

 	Start() error

	Stop() error

	Destroy() error
}
