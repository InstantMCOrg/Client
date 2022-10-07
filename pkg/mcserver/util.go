package mcserver

import "log"

func IsRunning() bool {
	return serverIsUpAndRunning
}

func SendStopCommand() {
	log.Println("Stopping Minecraft Server...")
	serverIsUpAndRunning = false
	SendCommand("stop")
}
