package mcserver

import (
	"fmt"
	"log"
)

func IsRunning() bool {
	return serverIsUpAndRunning
}

func SendStopCommand() {
	log.Println("Stopping Minecraft Server...")
	serverIsUpAndRunning = false
	SendCommand("stop")
}

func OpPlayer(targetPlayer string) {
	SendCommand(fmt.Sprintf("/op %s", targetPlayer))
}

func RamSize() int {
	return targetRamSize
}

func SetRamSize(ramSize int) {
	targetRamSize = ramSize
}
