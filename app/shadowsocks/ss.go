package shadowsocks


import (
	ss "github.com/JohnSmithX/mus/app/shadowsocks/lib"
	"time"
)


func New(addr, method, encrStr string, timeout time.Duration, fn func(*int))(server *ss.ProxyServer, err error) {
	config := &ss.ProxyConfig{
		Addr: addr,
		Method: method,
		EncrStr: encrStr,
		Timeout: timeout,
	}
	server, err = ss.New(config)

	server.CallbackMethods = &traffic{
		recordFunc: fn,
	}

	return
}
