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
