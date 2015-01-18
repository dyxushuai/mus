package models

type IServer interface {
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

	//db
	Update() error
	Delete() error
}
