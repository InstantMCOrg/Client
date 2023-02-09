package main

import (
	"github.com/instantminecraft/client/pkg/auth"
	"github.com/instantminecraft/client/pkg/constants"
	"github.com/instantminecraft/client/pkg/mcserver"
	"github.com/instantminecraft/client/pkg/proxy"
	"github.com/instantminecraft/client/pkg/router"
	"github.com/instantminecraft/client/pkg/server"
	"log"
	"os"
)

func main() {
	if !auth.HasAuthKey() {
		log.Println("Warning: \"auth\" environment variable has not been set. The http server will start without authentication.")
	}
	routes := router.Register()
	go server.Handle(routes)
	buildWorld, _ := os.LookupEnv(constants.EnvBuildMcWorldOnBoot)
	if buildWorld == "true" {
		mcserver.StartServer()
	}
	proxy.Start()
}
