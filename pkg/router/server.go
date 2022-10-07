package router

import (
	"github.com/instantminecraft/client/pkg/mcserver"
	"github.com/instantminecraft/client/pkg/server"
	"net/http"
)

func start(w http.ResponseWriter, r *http.Request) {
	if mcserver.IsRunning() {
		server.CreateResponse(w, "Minecraft Server already running", http.StatusConflict)
		return
	}

	blocking := r.URL.Query().Get("blocking") == "true"

	mcserver.StartServer()

	if blocking {
		mcserver.WaitUntilServerIsReady()
		server.CreateResponse(w, "Minecraft Server is running", http.StatusOK)
	} else {
		server.CreateResponse(w, "Minecraft Server has been started", http.StatusOK)
	}
}

func stop(w http.ResponseWriter, r *http.Request) {
	if !mcserver.IsRunning() {
		server.CreateResponse(w, "Can't stop an already stopped Minecraft Server", http.StatusConflict)
		return
	}

	blocking := r.URL.Query().Get("blocking") == "true"

	mcserver.SendStopCommand()

	if blocking {
		mcserver.WaitForStop()
		server.CreateResponse(w, "Minecraft Server has stopped", http.StatusOK)
	} else {
		server.CreateResponse(w, "Minecraft Server is stopping", http.StatusOK)
	}
}
