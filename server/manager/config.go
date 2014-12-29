
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
	//server command channel size
	serverCommand int = 10

)

const (
	localFormat string = "client: %s : %%s %%v"
)
