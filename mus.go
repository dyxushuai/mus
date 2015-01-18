package main

import (
	"github.com/JohnSmithX/mus/app"
	"github.com/JohnSmithX/mus/config"

)

func main() {
	app.Serve(config.REDIS_SERVER, config.REDIS_PASSWORD)
}
