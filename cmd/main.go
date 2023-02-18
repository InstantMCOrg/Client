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
	"strconv"
)

func main() {
	if !auth.HasAuthKey() {
		log.Println("Warning: \"auth\" environment variable has not been set. The http server will start without authentication.")
	}
	routes := router.Register()
	go server.Handle(routes)
	targetRamSizeRaw, ok := os.LookupEnv(constants.EnvTargetRamSize)
	if ok {
		targetRamSize, err := strconv.Atoi(targetRamSizeRaw)
		if err != nil {
			log.Fatal("Couldn't parse requested ram size:", err)
		} else if targetRamSize < constants.MinimumRamMb {
			log.Printf("Requested ram size (%dmb) can't be below minimum ram size (%dmb). Defaulting back to minimum ram size...\n", targetRamSize, constants.MinimumRamMb)
			targetRamSize = constants.MinimumRamMb
		}
		mcserver.SetRamSize(targetRamSize)
		log.Printf("Requested ram size has been set to %dmb\n", targetRamSize)
	}
	buildWorld, _ := os.LookupEnv(constants.EnvBuildMcWorldOnBoot)
	if buildWorld == "true" {
		mcserver.StartServer(mcserver.RamSize()) // using default or env value
	}
	proxy.Start()
}
