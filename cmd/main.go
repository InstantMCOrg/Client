package main

import (
	"github.com/instantminecraft/client/pkg/auth"
	"github.com/instantminecraft/client/pkg/proxy"
	"github.com/instantminecraft/client/pkg/router"
	"github.com/instantminecraft/client/pkg/server"
	"log"
)

func main() {
	if !auth.HasAuthKey() {
		log.Println("Warning: \"auth\" environment variable has not been set. The http server will start without authentication.")
	}
	routes := router.Register()
	go server.Handle(routes)
	proxy.Start()
}
