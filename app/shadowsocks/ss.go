package shadowsocks


import (
	ss "github.com/JohnSmithX/mus/app/shadowsocks/lib"
	"time"
)


type Proxyer interface {
	Listen()
	IsStopped() bool
	Stop()
	SetCallbacks(ss.CallbackInterface)
}

func New(addr, method, encrStr string, timeout time.Duration, fn func(*int))(server Proxyer, err error) {
	config := &ss.ProxyConfig{
		Addr: addr,
		Method: method,
		EncrStr: encrStr,
		Timeout: timeout,
	}
	server, err = ss.New(config)

	server.SetCallbacks(&traffic{
		recordFunc: fn,
	})

	return
}
