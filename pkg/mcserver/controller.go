package mcserver

import (
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

var isGeneratingWorld = false
var worldGeneratingStatus = 0
var WorldGenerationChan = make(chan int)
var serverIsUpAndRunning = false
var stdin io.WriteCloser
var cmd *exec.Cmd
var serverStartupLock sync.WaitGroup

const (
	PORT = 25565
)

func StartServer() {
	cmd = exec.Command("java", "-jar", "server.jar")
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		log.Fatal("Couldn't redirect output", err)
	}

	stdin, err = cmd.StdinPipe()

	if err != nil {
		log.Fatal("Couldn't redirect input", err)
	}

	serverStartupLock.Add(1)
	err = cmd.Start()

	if err != nil {
		log.Fatal("Couldn't start Minecraft Server", err)
	}
	log.Println("Started Minecraft Server successfully. Waiting...")

	go readServerOutput(stdout)
}

func SendCommand(command string) {
	if _, err := io.WriteString(stdin, command+"\n\r"); err != nil {
		log.Println("Couldn't write to Minecraft command pipe:", err)
	}
}

func WaitForStop() {
	cmd.Wait()
	serverIsUpAndRunning = false
	log.Println("Minecraft Server stopped")
}

func WaitUntilServerIsReady() {
	serverStartupLock.Wait()
}

func readServerOutput(pipe io.ReadCloser) {
	for {
		tmp := make([]byte, 1024)
		_, err := pipe.Read(tmp)
		parseMinecraftLog(tmp)
		if err != nil {
			break
		}
	}
}

func parseMinecraftLog(output []byte) {
	minecraftLog := string(output)

	if strings.Contains(minecraftLog, "Preparing spawn area: ") {
		spawnAreaPreparingPercent, err := strconv.Atoi(strings.Split(strings.Split(minecraftLog, "area: ")[1], "%")[0])
		log.Println("Spawning area is currently being prepared at", spawnAreaPreparingPercent, "%")
		isGeneratingWorld = true
		if err != nil {
			worldGeneratingStatus = spawnAreaPreparingPercent
		}
		select {
		case WorldGenerationChan <- spawnAreaPreparingPercent:
			// A websocket is currently listening...
			break
		default:
			// No websocket is currently listening
			break
		}
	} else if strings.Contains(minecraftLog, "Done") {
		isGeneratingWorld = false
		serverIsUpAndRunning = true
		log.Println("Minecraft Server is up and running")
		WorldGenerationChan <- 100
		serverStartupLock.Done()
	}
}
