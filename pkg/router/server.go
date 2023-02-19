package router

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/instantminecraft/client/pkg/mcserver"
	"github.com/instantminecraft/client/pkg/server"
	"log"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func start(w http.ResponseWriter, r *http.Request) {
	if mcserver.IsRunning() {
		server.CreateResponse(w, "Minecraft Server already running", http.StatusConflict)
		return
	}

	blocking := r.URL.Query().Get("blocking") == "true"
	ramMBRaw := r.URL.Query().Get("ram")
	setRam := ramMBRaw != ""
	if setRam {
		ramMB, err := strconv.Atoi(ramMBRaw)
		if err != nil {
			server.CreateResponse(w, "Couldn't parse \"ram\" parameter", http.StatusBadRequest)
			return
		}
		mcserver.StartServer(ramMB)
	} else {
		mcserver.StartServer(mcserver.RamSize()) // using already used ram size
	}

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

func creationStatus(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		server.CreateResponse(w, "Couldn't establish a websocket connection", http.StatusOK)
		return
	}
	if mcserver.IsRunning() {
		conn.WriteJSON(map[string]interface{}{"status": "already running"})
		conn.Close()
		return
	}

	for {
		currentGenerationStatus := <-mcserver.WorldGenerationChan
		err := conn.WriteJSON(map[string]interface{}{"status": "preparing", "world_status": currentGenerationStatus})
		if err != nil {
			log.Println("An error occurred while sending world generation status:", err)
			break
		}
		if currentGenerationStatus == 100 {
			break
		}
	}

	conn.Close()
}

func serverLogs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		server.CreateResponse(w, "Couldn't establish a websocket connection", http.StatusOK)
		return
	}

	for {
		message := <-mcserver.ServerLogsChan
		conn.WriteMessage(websocket.TextMessage, []byte(message))
	}

	conn.Close()
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	message := r.FormValue("message")
	if message == "" {
		server.CreateResponse(w, "Please provide \"message\"", http.StatusBadRequest)
		return
	}
	mcserver.SendCommand(fmt.Sprintf("/say %s", message))
	server.CreateResponse(w, "Success", http.StatusOK)
}
