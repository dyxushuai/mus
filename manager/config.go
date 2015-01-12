//default config
package manager

//broadcast channel size
const (
	numOfErr int = 1
	numOfMsg int = 100
)

//manager config
const (
)

//server config
const (
	//format output
	serverFormat string = "proxy server at port %s : %%s %%v"
)


//local config
const (
	localFormat string = "local: %s : %%s %%v"
)

//redis config
const (
	serverPrefix = "server:"
	flowPrefix = "flow:"
)
