package controllers

import (
	"github.com/JohnSmithX/mus/manager"
	"github.com/JohnSmithX/mus/config"
)

var M *manager.Manager

func init() {
	M = manager.New(config.REDIS_SERVER, config.REDIS_PASSWORD, config.VERBOSE)
}
