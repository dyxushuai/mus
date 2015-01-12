package controllers

import (
	"github.com/JohnSmithX/mus/manager"
	"github.com/JohnSmithX/mus/config"
	"fmt"
)

var M *manager.Manager

func init() {
	var err error
	M,err = manager.New(config.REDIS_SERVER, config.REDIS_PASSWORD, config.VERBOSE)
	fmt.Println(err)
}

func NewServerAPI() *ServerAPI {
	return &ServerAPI{}
}
