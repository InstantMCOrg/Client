package mcserver

import (
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

var isGeneratingWorld = true
var worldGeneratingStatus = 0
var serverIsUpAndRunning = false
var stdin io.WriteCloser
var cmd *exec.Cmd
var serverStartupLock sync.WaitGroup

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
	log.Println("Started Minecraft Server successfully")

	go readServerOutput(stdout)
}

func SendCommand(command string) {
	if _, err := io.WriteString(stdin, command+"\n\r"); err != nil {
		log.Println("Couldn't write to Minecraft command pipe:", err)
	}
}

func WaitForStop() {
	cmd.Wait()
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

	log.Println(minecraftLog)

	if strings.Contains(minecraftLog, "Preparing spawn area: ") {
		spawnAreaPreparingPercent, err := strconv.Atoi(strings.Split(strings.Split(minecraftLog, "area: ")[1], "%")[0])
		log.Println("Spawning area is currently being prepared at ", spawnAreaPreparingPercent, "%")
		isGeneratingWorld = true
		if err != nil {
			worldGeneratingStatus = spawnAreaPreparingPercent
		}
	} else if strings.Contains(minecraftLog, "Done") {
		isGeneratingWorld = false
		serverIsUpAndRunning = true
		log.Println("Minecraft Server is up and running")
		serverStartupLock.Done()
	}
}
